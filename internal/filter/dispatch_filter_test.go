package filter

import (
	"io"
	"net"
	"testing"

	"github.com/frizz925/higuchi/internal/dispatcher"
	"github.com/frizz925/higuchi/internal/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDispatchFilter(t *testing.T) {
	require := require.New(t)
	expected := []byte("expected")
	c1, c2 := net.Pipe()

	df := NewDispatchFilter(dispatcher.DispatcherFunc(func(rw io.ReadWriter, addr string) error {
		return testutil.EchoReadWriter(rw)
	}))
	errCh := make(chan error, 1)
	go func() {
		errCh <- df.Do(&Context{
			Conn:   c2,
			Logger: zap.NewExample(),
		}, "")
		close(errCh)
	}()

	buf := make([]byte, 512)
	_, err := c1.Write(expected)
	require.NoError(err)
	n, err := c1.Read(buf)
	require.NoError(err)
	require.NoError(<-errCh)
	require.Equal(expected, buf[:n])
}
