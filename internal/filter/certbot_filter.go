package filter

import (
	"net/http"
	"strings"

	"github.com/frizz925/higuchi/internal/httputil"
)

const AcmeChallengePath = "/.well-known/acme-challenge"

var AcmeChallengePathLength = len(AcmeChallengePath)

type CertbotFilter struct {
	hostname string
	handler  http.Handler
}

func NewCertbotFilter(hostname, webroot string) *CertbotFilter {
	return &CertbotFilter{
		hostname: hostname,
		handler:  http.FileServer(http.Dir(webroot)),
	}
}

func (cf *CertbotFilter) Do(ctx *Context, req *http.Request, next Next) error {
	if !cf.checkRequest(req) {
		return next()
	}
	req.URL.Path = req.URL.Path[AcmeChallengePathLength+1:]
	cw := httputil.NewResponseWriter(ctx, req)
	cf.handler.ServeHTTP(cw, req)
	return cw.Close()
}

func (cf *CertbotFilter) checkRequest(req *http.Request) bool {
	return req.Method == http.MethodGet && req.Host == cf.hostname && strings.HasPrefix(req.URL.Path, AcmeChallengePath)
}
