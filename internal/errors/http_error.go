package errors

import (
	"net"
	"net/url"
)

type HTTPError struct {
	Err         string
	User        *url.Userinfo
	Source      net.Addr
	Listener    net.Addr
	Destination string
	StatusCode  int
}

func (e *HTTPError) Error() string {
	return e.Err
}
