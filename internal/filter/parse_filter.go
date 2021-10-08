package filter

import (
	"bufio"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/frizz925/higuchi/internal/netutil"
)

const DefaultBufferSize = 512

type ParseFilter struct {
	filters []HTTPFilter
}

func NewParseFilter(filters ...HTTPFilter) *ParseFilter {
	return &ParseFilter{
		filters: filters,
	}
}

func (pf *ParseFilter) Do(ctx *Context, next Next) error {
	req, err := http.ReadRequest(bufio.NewReaderSize(ctx, DefaultBufferSize))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	req.URL.Scheme = "http"
	req.URL.Host = req.Host

	oldHeader := req.Header
	newHeader := make(http.Header)
	for k, v := range oldHeader {
		if strings.HasPrefix(k, "Proxy") {
			continue
		}
		newHeader[k] = v
	}
	req.Header = newHeader

	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return err
	}
	req.Header = oldHeader
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
