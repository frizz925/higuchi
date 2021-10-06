package filter

import (
	"net"
	"net/http"
)

type ForwardFilter struct {
	filters []NetFilter
}

func NewForwardFilter(filters ...NetFilter) *ForwardFilter {
	return &ForwardFilter{filters}
}

func (ff *ForwardFilter) Do(conn net.Conn, req *http.Request) error {
	hostport := req.Host
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		host = hostport
		port = "80"
	}
	addr := net.JoinHostPort(host, port)
	for _, f := range ff.filters {
		if err := f.Do(conn, addr); err != nil {
			return err
		}
	}
	return nil
}
