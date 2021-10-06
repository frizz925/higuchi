package dispatcher

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTCPDispatcher(t *testing.T) {
	require := require.New(t)
	expected := []byte("expected")

	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(err)
	laddr := l.Addr().String()

	sentCh := make(chan []byte, 1)
	errCh := make(chan error, 1)
	go func(errCh chan<- error) {
		defer close(errCh)
		c, err := l.Accept()
		if err != nil {
			errCh <- err
			return
		}
		sent := make([]byte, 512)
		n, err := c.Read(sent)
		if err != nil {
			errCh <- err
			return
		}
		sentCh <- sent[:n]
		_, err = c.Write(expected)
		if err != nil {
			errCh <- err
		}
	}(errCh)

	c1, c2 := net.Pipe()
	go DefaultTCPDispatcher.Dispatch(c2, laddr)
	_, err = c1.Write(expected)
	require.NoError(err)

	buf := make([]byte, 512)
	n, err := c1.Read(buf)
	require.NoError(err)
	require.NoError(<-errCh)
	recv := buf[:n]

	require.Equal(expected, recv)
	require.Equal(recv, <-sentCh)
}
