package httputil

import (
	"bufio"
	"bytes"
	"net/http"
	"testing"

	"github.com/frizz925/higuchi/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestParseRequestHeader(t *testing.T) {
	testTable := []struct {
		name          string
		rawBody       []byte
		assertFunc    func(assert *assert.Assertions, req *http.Request)
		expectedError error
	}{
		{
			name: "plain proxy request",
			rawBody: testutil.LinesToRawHeader(
				"GET http://example.org/ HTTP/1.1",
				"Host: example.org",
				"User-Agent: curl/7.64.1",
				"Accept: */*",
				"Proxy-Connection: Keep-Alive",
			),
			assertFunc: func(assert *assert.Assertions, req *http.Request) {
				assert.Equal(http.MethodGet, req.Method)
				assert.Equal("HTTP/1.1", req.Proto)
				assert.Equal("example.org", req.Host)
				assert.Equal("http://example.org/", req.URL.String())
			},
		},
		{
			name: "tunnel proxy request",
			rawBody: testutil.LinesToRawHeader(
				"CONNECT example.org:80 HTTP/1.1",
				"Host: example.org:80",
				"User-Agent: curl/7.64.1",
				"Proxy-Connection: Keep-Alive",
			),
			assertFunc: func(assert *assert.Assertions, req *http.Request) {
				assert.Equal(http.MethodConnect, req.Method)
				assert.Equal("HTTP/1.1", req.Proto)
				assert.Equal("example.org:80", req.Host)
				assert.Equal("example.org:80", req.URL.Host)
				assert.Equal("Keep-Alive", req.Header.Get("Proxy-Connection"))
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			br := bufio.NewReader(bytes.NewReader(tt.rawBody))
			req, err := ParseRequestHeader(br)
			if tt.expectedError != nil {
				assert.Equal(tt.expectedError, err)
			} else if err == nil && tt.assertFunc != nil {
				tt.assertFunc(assert, req)
			} else {
				assert.NoError(err)
			}
		})
	}
}
