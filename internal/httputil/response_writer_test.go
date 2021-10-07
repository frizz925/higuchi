package httputil

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResponseWriter(t *testing.T) {
	expected := []byte("expected")
	expectedLen := len(expected)

	buf := &bytes.Buffer{}
	req := &http.Request{
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Method:     http.MethodGet,
	}

	testTable := []struct {
		name    string
		chunked bool
	}{
		{
			name:    "fixed content length",
			chunked: false,
		},
		{
			name:    "chunked content",
			chunked: true,
		},
	}
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			rw := NewResponseWriter(buf, req)
			if !tt.chunked {
				rw.Header().Set("Content-Length", fmt.Sprint(expectedLen))
			}
			_, err := rw.Write(expected)
			require.NoError(err)
			require.NoError(rw.Close())

			res, err := http.ReadResponse(bufio.NewReader(buf), req)
			require.NoError(err)
			defer res.Body.Close()

			require.Equal(http.StatusOK, res.StatusCode)
			if !tt.chunked {
				require.Equal(int64(expectedLen), res.ContentLength)
			}
			b, err := ioutil.ReadAll(res.Body)
			require.NoError(err)
			require.Equal(expected, b)
		})
	}
}
