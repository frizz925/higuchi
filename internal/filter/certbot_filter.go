package filter

import (
	"net"
	"net/http"
	"strings"

	"github.com/frizz925/higuchi/internal/httputil"
	"go.uber.org/zap"
)

type CertbotFilter struct {
	hostname      string
	challengePath string
	handler       http.Handler
}

type CertbotConfig struct {
	Hostname      string
	Webroot       string
	ChallengePath string
}

func NewCertbotFilter(cfg CertbotConfig) *CertbotFilter {
	return &CertbotFilter{
		hostname:      cfg.Hostname,
		challengePath: cfg.ChallengePath,
		handler:       http.FileServer(http.Dir(cfg.Webroot)),
	}
}

func (cf *CertbotFilter) Do(ctx *Context, req *http.Request, next Next) error {
	if !cf.checkRequest(req) {
		return next()
	}
	ctx.Logger.Info("Captured certbot challenge request", zap.String("path", req.URL.Path))
	cw := httputil.NewResponseWriter(ctx, req)
	cf.handler.ServeHTTP(cw, req)
	return cw.Close()
}

func (cf *CertbotFilter) checkRequest(req *http.Request) bool {
	hostport := req.Host
	host, _, err := net.SplitHostPort(hostport)
	if err != nil {
		host = hostport
	}
	return req.Method == http.MethodGet && host == cf.hostname && strings.HasPrefix(req.URL.Path, cf.challengePath)
}
