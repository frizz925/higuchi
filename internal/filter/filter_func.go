package filter

import "net/http"

type FilterFunc func(c *Context, next Next) error

func (fn FilterFunc) Do(c *Context, next Next) error {
	return fn(c, next)
}

type HTTPFilterFunc func(c *Context, req *http.Request, next Next) error

func (fn HTTPFilterFunc) Do(c *Context, req *http.Request, next Next) error {
	return fn(c, req, next)
}

type NetFilterFunc func(c *Context, addr string, next Next) error

func (fn NetFilterFunc) Do(c *Context, addr string, next Next) error {
	return fn(c, addr, next)
}
