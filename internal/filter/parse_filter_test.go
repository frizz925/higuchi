package filter

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/frizz925/higuchi/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestParseFilter(t *testing.T) {
	require := require.New(t)
	expectedBody, expectedHost := "expected", "localhost"
	pf := NewParseFilter(HTTPFilterFunc(func(conn net.Conn, req *http.Request) error {
		defer req.Body.Close()
		var buf bytes.Buffer
		res := fmt.Sprintf("host=%s body=", req.Host)
		if _, err := buf.WriteString(res); err != nil {
			return err
		}
		if _, err := buf.ReadFrom(req.Body); err != nil && err != io.EOF {
			return err
		}
		if _, err := buf.WriteTo(conn); err != nil {
			return err
		}
		return nil
	}))

	c1, c2 := net.Pipe()
	errCh := make(chan error, 1)
	go func() {
		errCh <- pf.Do(c2)
		close(errCh)
	}()

	_, err := c1.Write(testutil.LinesToRawRequest(
		fmt.Sprintf("GET http://%s/ HTTP/1.1", expectedHost),
		fmt.Sprintf("Host: %s", expectedHost),
		"User-Agent: curl/7.64.1",
		"Accept: */*",
		"Proxy-Connection: Keep-Alive",
		"Content-Type: text/plain",
		fmt.Sprintf("Content-Length: %d", len(expectedBody)),
		"",
		expectedBody,
	))
	require.NoError(err)

	buf := make([]byte, 512)
	n, err := c1.Read(buf)
	require.NoError(err)

	var (
		host string
		body string
	)
	s := string(buf[:n])
	_, err = fmt.Sscanf(s, "host=%s body=%s", &host, &body)
	require.NoError(err)
	require.Equal(expectedHost, host)
	require.Equal(expectedBody, body)
}
