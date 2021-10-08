package dispatcher

import (
	"io"
	"net"
	"sync"

	"github.com/frizz925/higuchi/internal/ioutil"
)

type TCPDispatcher struct {
	ioPool sync.Pool
}

func NewTCPDispatcher(bufsize int) *TCPDispatcher {
	d := &TCPDispatcher{}
	d.ioPool.New = func() interface{} {
		return &ioBuffer{
			sbuf: make([]byte, bufsize),
			dbuf: make([]byte, bufsize),
		}
	}
	return d
}

func (d *TCPDispatcher) Dispatch(rw io.ReadWriter, addr string) error {
	out, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer out.Close()
	iob := d.ioPool.Get().(*ioBuffer)
	defer d.ioPool.Put(iob)
	return ioutil.PipeBuffer(rw, out, iob.sbuf, iob.dbuf)
}
