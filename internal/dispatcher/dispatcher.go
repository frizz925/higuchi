package dispatcher

import "io"

type Dispatcher interface {
	Dispatch(rw io.ReadWriter, addr string) error
}
