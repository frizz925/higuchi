package filter

import (
	"net"
	"net/http"

	"github.com/frizz925/higuchi/internal/errors"
	"go.uber.org/zap"
)

type Context struct {
	net.Conn
	Logger *zap.Logger
}

type Filter interface {
	Do(c *Context) error
}

type HTTPFilter interface {
	Do(c *Context, req *http.Request) error
}

type NetFilter interface {
	Do(c *Context, addr string) error
}

func ToHTTPError(ctx *Context, req *http.Request, err string, code int) *errors.HTTPError {
	e := &errors.HTTPError{
		Err:         err,
		Source:      ctx.RemoteAddr(),
		Listener:    ctx.LocalAddr(),
		Destination: req.Host,
		Request:     req,
		StatusCode:  code,
	}
	if req.URL != nil && req.URL.User != nil {
		e.User = req.URL.User
	}
	return e
}
