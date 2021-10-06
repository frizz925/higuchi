package ioutil

import (
	"io"

	"golang.org/x/sync/errgroup"
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
	var g errgroup.Group
	g.Go(func() error {
		return pipe(src, dst, sbuf)
	})
	g.Go(func() error {
		return pipe(dst, src, dbuf)
	})
	return g.Wait()
}

func pipe(src io.Reader, dst io.Writer, buf []byte) error {
	_, err := io.CopyBuffer(dst, src, buf)
	return err
}
