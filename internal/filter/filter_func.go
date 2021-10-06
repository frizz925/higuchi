package filter

import "net/http"

type FilterFunc func(c *Context) error

func (fn FilterFunc) Do(c *Context) error {
	return fn(c)
}

type HTTPFilterFunc func(c *Context, req *http.Request) error

func (fn HTTPFilterFunc) Do(c *Context, req *http.Request) error {
	return fn(c, req)
}

type NetFilterFunc func(c *Context, addr string) error

func (fn NetFilterFunc) Do(c *Context, addr string) error {
	return fn(c, addr)
}
