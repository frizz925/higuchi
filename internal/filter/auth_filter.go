package filter

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

type AuthCompareFunc func(password string, v interface{}) bool

var AuthCompareString = func(password string, v interface{}) bool {
	switch v.(type) {
	case string:
		return password == v
	default:
		return false
	}
}

type AuthFilter struct {
	users   map[string]interface{}
	compare AuthCompareFunc
}

func NewAuthFilter(users map[string]interface{}, compare ...AuthCompareFunc) *AuthFilter {
	cmp := AuthCompareString
	return &AuthFilter{users, cmp}
}

func (af *AuthFilter) Do(c *Context, req *http.Request, next Next) error {
	auth := req.Header.Get("Proxy-Authorization")
	if auth == "" {
		return ToHTTPError(c, req, "authorization required", http.StatusProxyAuthRequired)
	}

	parts := strings.Split(auth, " ")
	if len(parts) < 2 {
		return ToHTTPError(c, req, "malformed authorization value", http.StatusBadRequest)
	}
	scheme, param := parts[0], parts[1]
	if scheme != "Basic" {
		return ToHTTPError(c, req, "unsupported authorization scheme", http.StatusBadRequest)
	}
	b, err := base64.StdEncoding.DecodeString(param)
	if err != nil {
		return ToHTTPError(c, req, "malformed authorization value", http.StatusBadRequest)
	}
	creds := strings.Split(string(b), ":")
	if len(creds) < 2 {
		return ToHTTPError(c, req, "malformed authorization param", http.StatusBadRequest)
	}

	user, pass := creds[0], creds[1]
	if req.URL == nil {
		req.URL = &url.URL{}
	}
	req.URL.User = url.UserPassword(user, pass)
	c.Logger = c.Logger.With(zap.String("user", user))

	v, ok := af.users[user]
	if !ok || !af.compare(pass, v) {
		return ToHTTPError(c, req, "invalid credentials", http.StatusForbidden)
	}
	return next()
}
