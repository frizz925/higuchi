package worker

import (
	"net"
	"testing"

	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWorker(t *testing.T) {
	require := require.New(t)
	expected := []byte("expected")

	w := New(filter.FilterFunc(func(c *filter.Context) error {
		return testutil.EchoReadWriter(c)
	}))

	c1, c2 := net.Pipe()
	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Handle(&filter.Context{
			Conn:   c2,
			Logger: zap.NewExample(),
		})
		close(errCh)
	}()

	_, err := c1.Write(expected)
	require.NoError(err)

	buf := make([]byte, 512)
	n, err := c1.Read(buf)
	require.Equal(expected, buf[:n])
}
