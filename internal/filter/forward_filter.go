package filter

import (
	"net/http"

	"github.com/frizz925/higuchi/internal/httputil"
)

type ForwardFilter struct {
	filters []NetFilter
}

func NewForwardFilter(filters ...NetFilter) *ForwardFilter {
	return &ForwardFilter{filters}
}

func (ff *ForwardFilter) Do(ctx *Context, req *http.Request, next Next) error {
	addr := httputil.ParseRequestAddress(req)
	ctx.LogFields.Destination = addr
	ctx.UpdateLogger()

	var netNext Next
	idx := 0
	netNext = func() error {
		if idx >= len(ff.filters) {
			return next()
		}
		f := ff.filters[idx]
		idx++
		return f.Do(ctx, addr, netNext)
	}
	if err := netNext(); err != nil {
		return ToHTTPError(ctx, req, err.Error(), http.StatusServiceUnavailable)
	}
	return nil
}
