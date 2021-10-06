package filter

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/frizz925/higuchi/internal/netutil"
)

const DefaultBufferSize = 512

type ParseFilter struct {
	filters []HTTPFilter
	buffer  *bufio.Reader
}

func NewParseFilter(filters ...HTTPFilter) *ParseFilter {
	return &ParseFilter{
		filters: filters,
		buffer:  bufio.NewReaderSize(nil, DefaultBufferSize),
	}
}

func (pf *ParseFilter) Do(conn net.Conn) error {
	pf.buffer.Reset(conn)
	req, err := http.ReadRequest(pf.buffer)
	if err != nil {
		return err
	}
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err
	}
	conn = netutil.NewPrefixedConn(conn, b)
	for _, f := range pf.filters {
		if err := f.Do(conn, req); err != nil {
			return err
		}
	}
	return nil
}
