package httputil

import (
	"bufio"
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteResponseHeader(t *testing.T) {
	require := require.New(t)
	buf := &bytes.Buffer{}
	expected := &http.Response{
		Proto:      "HTTP/1.1",
		StatusCode: http.StatusOK,
		Header: http.Header{
			"X-Expected-Header": []string{"expected"},
		},
	}

	bw := bufio.NewWriter(buf)
	WriteResponseHeader(expected, bw)
	require.NoError(bw.Flush())

	res, err := http.ReadResponse(bufio.NewReader(buf), nil)
	require.NoError(err)
	require.Equal(expected.Proto, res.Proto)
	require.Equal(expected.StatusCode, res.StatusCode)
	require.Equal(expected.Header.Get("X-Expected-Header"), res.Header.Get("X-Expected-Header"))
}
