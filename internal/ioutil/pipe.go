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
	var g errgroup.Group
	g.Go(func() error {
		return pipe(src, dst, n)
	})
	g.Go(func() error {
		return pipe(dst, src, n)
	})
	return g.Wait()
}

func pipe(src io.Reader, dst io.Writer, n int) error {
	buf := make([]byte, n)
	_, err := io.CopyBuffer(dst, src, buf)
	return err
}
