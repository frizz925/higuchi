package filter

import (
	"net"
	"net/http"

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
