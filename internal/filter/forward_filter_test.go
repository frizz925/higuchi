package filter

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestForwardFilter(t *testing.T) {
	require := require.New(t)
	ff := NewForwardFilter(NetFilterFunc(func(c *Context, addr string) error {
		_, err := c.Write([]byte(addr))
		return err
	}))

	c1, c2 := net.Pipe()
	errCh := make(chan error, 1)
	go func() {
		errCh <- ff.Do(&Context{
			Conn:   c2,
			Logger: zap.NewExample(),
		}, &http.Request{
			Method: http.MethodGet,
			Host:   "localhost",
		})
		close(errCh)
	}()

	buf := make([]byte, 64)
	n, err := c1.Read(buf)
	require.NoError(err)
	require.Equal("localhost:80", string(buf[:n]))
}
