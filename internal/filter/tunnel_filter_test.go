package filter

import (
	"bufio"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/frizz925/higuchi/internal/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTunnelFilter(t *testing.T) {
	expected := []byte("expected")
	require := require.New(t)
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(err)

	errCh := make(chan error, 1)
	go func() {
		defer l.Close()
		c, err := l.Accept()
		if err != nil {
			errCh <- err
			return
		}
		defer c.Close()
		errCh <- testutil.EchoReadWriter(c)
	}()

	c1, c2 := net.Pipe()
	req := &http.Request{
		Method: http.MethodConnect,
		Host:   l.Addr().String(),
		URL: &url.URL{
			Host: l.Addr().String(),
		},
		Proto: "HTTP/1.1",
	}
	go NewTunnelFilter(512).Do(&Context{
		Conn:   c2,
		Logger: zap.NewExample(),
	}, req, nil)

	res, err := http.ReadResponse(bufio.NewReader(c1), req)
	require.NoError(err)
	require.Equal(http.StatusOK, res.StatusCode)
	require.Equal("200 Connection established", res.Status)

	_, err = c1.Write(expected)
	require.NoError(err)

	buf := make([]byte, 512)
	n, err := c1.Read(buf)
	require.NoError(<-errCh)
	require.NoError(err)
	require.Equal(expected, buf[:n])
}
