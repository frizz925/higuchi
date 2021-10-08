package filter

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHealthCheckFilter(t *testing.T) {
	require := require.New(t)
	f := NewHealthCheckFilter(http.MethodOptions, "/")

	c1, c2 := net.Pipe()
	errCh := make(chan error, 1)
	go func() {
		defer c2.Close()
		errCh <- f.Do(
			NewContext(c2, zap.NewExample()),
			NextError(errors.New("shouldn't call next")),
		)
		close(errCh)
	}()

	brw := bufio.NewReadWriter(bufio.NewReader(c1), bufio.NewWriter(c1))
	brw.WriteString("OPTIONS / HTTP/1.1\r\n\r\n")
	require.NoError(brw.Flush())

	res, err := http.ReadResponse(bufio.NewReader(c1), nil)
	require.NoError(<-errCh)
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode)
}
