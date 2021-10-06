package dispatcher

import (
	"io"
	"net"

	"github.com/frizz925/higuchi/internal/ioutil"
)

var DefaultTCPDispatcher = &TCPDispatcher{}

type TCPDispatcher struct {
	net.Dialer
}

func (d *TCPDispatcher) Dispatch(rw io.ReadWriter, addr string) error {
	out, err := d.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer out.Close()
	return ioutil.Pipe(rw, out)
}
