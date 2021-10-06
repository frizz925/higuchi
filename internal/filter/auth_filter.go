package filter

import (
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"strings"
)

type AuthFilter struct {
	users map[string]string
}

func NewAuthFilter(users map[string]string) *AuthFilter {
	return &AuthFilter{users}
}

// TODO: Create standardized errors
func (af *AuthFilter) Do(conn net.Conn, req *http.Request) error {
	auth := req.Header.Get("Proxy-Authorization")
	if auth == "" {
		return errors.New("authorization required")
	}
	parts := strings.Split(auth, " ")
	if len(parts) < 2 {
		return errors.New("malformed authorization value")
	}
	scheme, param := parts[0], parts[1]
	if scheme != "Basic" {
		return errors.New("unsupported authorization scheme")
	}
	b, err := base64.StdEncoding.DecodeString(param)
	if err != nil {
		return errors.New("failed to decode authorization param")
	}
	creds := strings.Split(string(b), ":")
	if len(creds) < 2 {
		return errors.New("malformed authorization param")
	}
	user, pass := creds[0], creds[1]
	v, ok := af.users[user]
	if !ok || pass != v {
		return errors.New("invalid credentials")
	}
	return nil
}
