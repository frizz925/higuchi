package filter

import (
	"net"
	"net/http"
)

type ForwardFilter struct {
	filters []NetFilter
}

func NewForwardFilter(filters ...NetFilter) *ForwardFilter {
	return &ForwardFilter{filters}
}

func (ff *ForwardFilter) Do(ctx *Context, req *http.Request, next Next) error {
	hostport := req.Host
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		host = hostport
		port = "80"
	}
	addr := net.JoinHostPort(host, port)
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
