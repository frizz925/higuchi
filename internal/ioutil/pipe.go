package ioutil

import (
	"io"
)

const DefaultPipeBufferSize = 1024

func Pipe(src, dst io.ReadWriter) error {
	return PipeSize(src, dst, DefaultPipeBufferSize)
}

func PipeSize(src, dst io.ReadWriter, n int) error {
	sbuf, dbuf := make([]byte, n), make([]byte, n)
	return PipeBuffer(src, dst, sbuf, dbuf)
}

func PipeBuffer(src, dst io.ReadWriter, sbuf, dbuf []byte) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- pipe(src, dst, sbuf)
	}()
	go func() {
		errCh <- pipe(dst, src, dbuf)
	}()
	return <-errCh
}

func pipe(src io.Reader, dst io.Writer, buf []byte) error {
	_, err := io.CopyBuffer(dst, src, buf)
	return err
}
