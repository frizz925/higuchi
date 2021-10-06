package filter

import (
	"net"

	"github.com/frizz925/higuchi/internal/dispatcher"
)

type DispatchFilter struct {
	dispatcher.Dispatcher
}

func NewDispatchFilter(d dispatcher.Dispatcher) *DispatchFilter {
	return &DispatchFilter{d}
}

func (df *DispatchFilter) Do(conn net.Conn, addr string) error {
	return df.Dispatch(conn, addr)
}
