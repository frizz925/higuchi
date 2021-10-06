package filter

import (
	"net"
	"net/http"
)

type Filter interface {
	Do(conn net.Conn) error
}

type HTTPFilter interface {
	Do(conn net.Conn, req *http.Request) error
}

type NetFilter interface {
	Do(conn net.Conn, addr string) error
}
