package filter

import (
	"net"
	"net/http"
)

type FilterFunc func(conn net.Conn) error

func (fn FilterFunc) Do(conn net.Conn) error {
	return fn(conn)
}

type HTTPFilterFunc func(conn net.Conn, req *http.Request) error

func (fn HTTPFilterFunc) Do(conn net.Conn, req *http.Request) error {
	return fn(conn, req)
}

type NetFilterFunc func(conn net.Conn, addr string) error

func (fn NetFilterFunc) Do(conn net.Conn, addr string) error {
	return fn(conn, addr)
}
