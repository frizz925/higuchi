package dispatcher

import "io"

type DispatcherFunc func(rw io.ReadWriter, addr string) error

func (d DispatcherFunc) Dispatch(rw io.ReadWriter, addr string) error {
	return d(rw, addr)
}
