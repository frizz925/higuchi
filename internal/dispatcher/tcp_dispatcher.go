package dispatcher

import (
	"io"
	"net"

	"github.com/frizz925/higuchi/internal/ioutil"
)

type TCPDispatcher struct {
	sbuf, dbuf []byte
}

func NewTCPDispatcher(bufsize int) *TCPDispatcher {
	return &TCPDispatcher{
		sbuf: make([]byte, bufsize),
		dbuf: make([]byte, bufsize),
	}
}

func (d *TCPDispatcher) Dispatch(rw io.ReadWriter, addr string) error {
	out, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer out.Close()
	return ioutil.PipeBuffer(rw, out, d.sbuf, d.dbuf)
}
