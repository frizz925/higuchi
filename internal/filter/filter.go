package filter

import (
	"net/http"

	"github.com/frizz925/higuchi/internal/errors"
)

var NextNoop Next = func() error {
	return nil
}

type Next func() error

type Filter interface {
	Do(c *Context, next Next) error
}

type HTTPFilter interface {
	Do(c *Context, req *http.Request, next Next) error
}

type NetFilter interface {
	Do(c *Context, addr string, next Next) error
}

func ToHTTPError(ctx *Context, req *http.Request, err string, code int) *errors.HTTPError {
	e := &errors.HTTPError{
		Err:         err,
		Source:      ctx.RemoteAddr(),
		Listener:    ctx.LocalAddr(),
		Destination: req.Host,
		Request:     req,
		StatusCode:  code,
		Header:      http.Header{},
	}
	if req.URL != nil && req.URL.User != nil {
		e.User = req.URL.User
	}
	return e
}

func NextError(err error) Next {
	return func() error {
		return err
	}
}
