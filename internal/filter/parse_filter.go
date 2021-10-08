package filter

import (
	"bufio"
	"net/http"
	"net/http/httputil"
	"strings"
	"sync"

	"github.com/frizz925/higuchi/internal/netutil"
)

const DefaultBufferSize = 512

type ParseFilter struct {
	filters []HTTPFilter
	bufPool sync.Pool
}

func NewParseFilter(bufsize int, filters ...HTTPFilter) *ParseFilter {
	pf := &ParseFilter{
		filters: filters,
	}
	pf.bufPool.New = func() interface{} {
		return bufio.NewReaderSize(nil, bufsize)
	}
	return pf
}

func (pf *ParseFilter) Do(ctx *Context, next Next) error {
	rd := pf.bufPool.Get().(*bufio.Reader)
	rd.Reset(ctx)

	req, err := http.ReadRequest(rd)
	if err != nil {
		pf.bufPool.Put(rd)
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

	pf.bufPool.Put(rd)
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
