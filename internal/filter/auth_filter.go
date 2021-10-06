package filter

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/frizz925/higuchi/internal/errors"
	"go.uber.org/zap"
)

var ()

type AuthFilter struct {
	users map[string]string
}

func NewAuthFilter(users map[string]string) *AuthFilter {
	return &AuthFilter{users}
}

// TODO: Create standardized errors
func (af *AuthFilter) Do(c *Context, req *http.Request) error {
	auth := req.Header.Get("Proxy-Authorization")
	if auth == "" {
		return errFromContextAndRequest(c, req, "authorization required", http.StatusProxyAuthRequired)
	}
	parts := strings.Split(auth, " ")
	if len(parts) < 2 {
		return errFromContextAndRequest(c, req, "malformed authorization value", http.StatusBadRequest)
	}
	scheme, param := parts[0], parts[1]
	if scheme != "Basic" {
		return errFromContextAndRequest(c, req, "unsupported authorization scheme", http.StatusBadRequest)
	}
	b, err := base64.StdEncoding.DecodeString(param)
	if err != nil {
		return errFromContextAndRequest(c, req, "malformed authorization value", http.StatusBadRequest)
	}
	creds := strings.Split(string(b), ":")
	if len(creds) < 2 {
		return errFromContextAndRequest(c, req, "malformed authorization param", http.StatusBadRequest)
	}

	user, pass := creds[0], creds[1]
	if req.URL == nil {
		req.URL = &url.URL{}
	}
	req.URL.User = url.UserPassword(user, pass)
	c.Logger = c.Logger.With(zap.String("user", user))

	v, ok := af.users[user]
	if !ok || pass != v {
		return errFromContextAndRequest(c, req, "invalid credentials", http.StatusForbidden)
	}
	return nil
}

func errFromContextAndRequest(c *Context, req *http.Request, err string, statusCode int) *errors.HTTPError {
	e := &errors.HTTPError{
		Err:         err,
		Source:      c.RemoteAddr(),
		Listener:    c.LocalAddr(),
		Destination: req.Host,
		StatusCode:  statusCode,
	}
	if req.URL != nil && req.URL.User != nil {
		e.User = req.URL.User
	}
	return e
}
