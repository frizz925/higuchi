package errors

import (
	"net"
	"net/url"
)

type EstablishedError struct {
	Err         string
	User        url.Userinfo
	Source      net.Addr
	Destination net.Addr
}

func (e *EstablishedError) Error() string {
	return e.Err
}
