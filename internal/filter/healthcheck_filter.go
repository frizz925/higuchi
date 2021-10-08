package filter

import (
	"bufio"
	"net/http"

	"github.com/frizz925/higuchi/internal/httputil"
)

type HealthCheckFilter struct {
	method, path string
}

func NewHealthCheckFilter(method, path string) *HealthCheckFilter {
	return &HealthCheckFilter{method, path}
}

func (hcf *HealthCheckFilter) Do(ctx *Context, req *http.Request, next Next) error {
	if req.Method != hcf.method || req.URL.Path != hcf.path {
		return next()
	}
	bw := bufio.NewWriter(ctx)
	httputil.WriteResponseHeader(&http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
		Request:    req,
	}, bw)
	return bw.Flush()
}
