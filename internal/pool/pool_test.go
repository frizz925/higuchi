package pool

import (
	"net"
	"testing"

	"github.com/frizz925/higuchi/internal/filter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func testPool(t *testing.T, p Pool) {
	require := require.New(t)
	expected := []byte("expected")
	c1, c2 := net.Pipe()
	errCh := make(chan error, 1)
	p.Dispatch(&filter.Context{
		Conn:   c2,
		Logger: zap.NewExample(),
	}, func(ctx *filter.Context, err error) {
		defer ctx.Close()
		errCh <- err
		close(errCh)
	})

	_, err := c1.Write(expected)
	require.NoError(err)
	buf := make([]byte, 512)
	n, err := c1.Read(buf)
	require.NoError(<-errCh)
	require.NoError(err)
	require.Equal(expected, buf[:n])
}
