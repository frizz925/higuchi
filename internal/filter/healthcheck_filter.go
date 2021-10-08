package filter

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/frizz925/higuchi/internal/httputil"
	"github.com/frizz925/higuchi/internal/netutil"
)

type HealthCheckFilter struct {
	method, path string
	br           *bufio.Reader
}

func NewHealthCheckFilter(method, path string) *HealthCheckFilter {
	return &HealthCheckFilter{
		method: method,
		path:   path,
		br:     bufio.NewReader(nil),
	}
}

func (hcf *HealthCheckFilter) Do(ctx *Context, next Next) error {
	hcf.br.Reset(ctx)
	b, _, err := hcf.br.ReadLine()
	if err != nil {
		return err
	}

	var method, path, proto string
	_, err = fmt.Sscanf(
		string(b), "%s %s %s",
		&method, &path, &proto,
	)
	if err != nil {
		return err
	}

	var protoMajor, protoMinor int
	_, err = fmt.Sscanf(proto, "HTTP/%d.%d", &protoMajor, &protoMinor)
	if err != nil {
		return err
	}

	if method != hcf.method || path != hcf.path {
		b = append(b, '\r', '\n')
		n := hcf.br.Buffered()
		if n > 0 {
			p := make([]byte, n)
			n, err := hcf.br.Read(p)
			if err != nil {
				return err
			}
			b = append(b, p[:n]...)
		}
		ctx.Conn = netutil.NewPrefixedConn(ctx.Conn, b)
		return next()
	}

	bw := bufio.NewWriter(ctx)
	httputil.WriteResponseHeader(&http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Proto:      proto,
		ProtoMajor: protoMajor,
		ProtoMinor: protoMinor,
	}, bw)
	if err := bw.Flush(); err != nil {
		return err
	}
	ctx.Logger.Info("Healthcheck called")
	return nil
}
