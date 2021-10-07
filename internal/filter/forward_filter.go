package filter

import (
	"net"
	"net/http"

	"github.com/frizz925/higuchi/internal/errors"
	"go.uber.org/zap"
)

type ForwardFilter struct {
	filters []NetFilter
}

func NewForwardFilter(filters ...NetFilter) *ForwardFilter {
	return &ForwardFilter{filters}
}

func (ff *ForwardFilter) Do(c *Context, req *http.Request) error {
	hostport := req.Host
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		host = hostport
		port = "80"
	}
	addr := net.JoinHostPort(host, port)
	c.Logger = c.Logger.With(zap.String("dst", addr))
	for _, f := range ff.filters {
		if err = f.Do(c, addr); err != nil {
			break
		}
	}
	if err != nil {
		he := &errors.HTTPError{
			Err:         err.Error(),
			Source:      c.RemoteAddr(),
			Listener:    c.LocalAddr(),
			Destination: addr,
			Request:     req,
			StatusCode:  http.StatusServiceUnavailable,
		}
		if req.URL != nil && req.URL.User != nil {
			he.User = req.URL.User
		}
		return he
	}
	return nil
}
