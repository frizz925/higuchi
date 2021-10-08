package filter

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHealthCheckFilter(t *testing.T) {
	require := require.New(t)
	f := NewHealthCheckFilter(http.MethodOptions, "*")
	req := &http.Request{
		Method: http.MethodOptions,
		URL: &url.URL{
			Path: "*",
		},
		Host: "localhost",
	}

	c1, c2 := net.Pipe()
	errCh := make(chan error, 1)
	go func() {
		defer c2.Close()
		errCh <- f.Do(
			NewContext(c2, zap.NewExample()), req,
			NextError(errors.New("shouldn't call next")),
		)
		close(errCh)
	}()

	res, err := http.ReadResponse(bufio.NewReader(c1), req)
	require.NoError(<-errCh)
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode)
}
