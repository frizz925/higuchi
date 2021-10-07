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

func (pf *ParseFilter) Do(ctx *Context, next Next) error {
	pf.buffer.Reset(ctx)
	req, err := http.ReadRequest(pf.buffer)
	if err != nil {
		return err
	}
	req.Header.Del("Proxy-Connection")
	req.URL.Scheme = "http"
	req.URL.Host = req.Host

	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err
	}
	ctx.Conn = netutil.NewPrefixedConn(ctx.Conn, b)

	var httpNext Next
	idx := 0
	httpNext = func() error {
		if idx >= len(pf.filters) {
			return next()
		}
		f := pf.filters[idx]
		idx++
		return f.Do(ctx, req, httpNext)
	}
	return httpNext()
}
