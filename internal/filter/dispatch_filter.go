package filter

import (
	"github.com/frizz925/higuchi/internal/dispatcher"
)

type DispatchFilter struct {
	dispatcher.Dispatcher
}

func NewDispatchFilter(d dispatcher.Dispatcher) *DispatchFilter {
	return &DispatchFilter{d}
}

func (df *DispatchFilter) Do(c *Context, addr string, _ Next) error {
	c.Logger.Info("Dispatching connection")
	return df.Dispatch(c, addr)
}
