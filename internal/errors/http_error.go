package errors

import (
	"net"
	"net/http"
	"net/url"
)

type HTTPError struct {
	Err         string
	User        *url.Userinfo
	Source      net.Addr
	Listener    net.Addr
	Destination string
	Request     *http.Request
	StatusCode  int
}

func (e *HTTPError) Error() string {
	return e.Err
}
