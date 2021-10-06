package filter

import (
	"bufio"
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

func (pf *ParseFilter) Do(c *Context) error {
	pf.buffer.Reset(c)
	req, err := http.ReadRequest(pf.buffer)
	if err != nil {
		return err
	}
	req.URL.Scheme = "http"
	req.URL.Host = req.Host
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err
	}
	c.Conn = netutil.NewPrefixedConn(c.Conn, b)
	for _, f := range pf.filters {
		if err := f.Do(c, req); err != nil {
			return err
		}
	}
	return nil
}
