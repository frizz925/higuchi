package netutil

import "net"

type PrefixedConn struct {
	net.Conn
	b []byte
}

func NewPrefixedConn(c net.Conn, b []byte) *PrefixedConn {
	return &PrefixedConn{c, b}
}

func (c *PrefixedConn) Read(b []byte) (int, error) {
	if len(c.b) <= 0 {
		return c.Conn.Read(b)
	}
	n := copy(b, c.b)
	c.b = c.b[n:]
	return n, nil
}
