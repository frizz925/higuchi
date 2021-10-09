package filter

import (
	"net/http"
)

var DefaultForwardedFilter = &ForwardedFilter{}

type ForwardedFilter struct{}

func (ff ForwardedFilter) Do(ctx *Context, req *http.Request, next Next) error {
	if proto := req.Header.Get("X-Forwarded-Proto"); proto != "" {
		ctx.LogFields.Proto = proto
	}
	if src := req.Header.Get("X-Forwarded-For"); src != "" {
		ctx.LogFields.Source = src
	}
	if host := req.Header.Get("X-Forwarded-Host"); host != "" {
		ctx.LogFields.Server = host
	}
	ctx.UpdateLogger()
	return next()
}
