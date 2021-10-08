package netutil

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrefixedConn(t *testing.T) {
	require := require.New(t)
	prefix, expected := []byte("prefix"), []byte("expected")

	c1, c2 := net.Pipe()
	go c1.Write(expected)

	pr := NewPrefixedConn(c2, prefix)
	buf := make([]byte, 512)

	var (
		n   int
		err error
	)

	n, err = pr.Read(buf)
	require.NoError(err)
	require.Equal(prefix, buf[:n])

	n, err = pr.Read(buf)
	require.NoError(err)
	require.Equal(expected, buf[:n])
}
