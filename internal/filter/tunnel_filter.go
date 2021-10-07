package filter

import (
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/frizz925/higuchi/internal/ioutil"
	"github.com/frizz925/higuchi/internal/netutil"
)

type TunnelFilter struct {
	sbuf, dbuf []byte

	fwTunnel *ForwardFilter
	tunCh    chan net.Conn
}

// Tunnel filter is a filter to intercept tunneling proxy request.
// The filter has its own chain when intercepting a tunneling connection.
// Otherwise, it would just use the provided chain instead.
func NewTunnelFilter(bufsize int) *TunnelFilter {
	tf := &TunnelFilter{
		sbuf:  make([]byte, bufsize),
		dbuf:  make([]byte, bufsize),
		tunCh: make(chan net.Conn, 1),
	}
	tf.fwTunnel = NewForwardFilter(NetFilterFunc(func(c *Context, addr string, _ Next) error {
		out, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		tf.tunCh <- out
		return nil
	}))
	return tf
}

func (tf *TunnelFilter) Do(ctx *Context, req *http.Request, next Next) error {
	if req.Method != http.MethodConnect {
		return next()
	}
	// Remove the prefixed buffer if we're tunneling
	if v, ok := ctx.Conn.(*netutil.PrefixedConn); ok {
		ctx.Conn = v.Conn
	}
	// Use forward filter to parse the address and grab the established connection from the channel
	if err := tf.fwTunnel.Do(ctx, req, nil); err != nil {
		return err
	}
	ctx.Logger.Info("Established connection for tunneling")
	b, err := httputil.DumpResponse(&http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 Connection established",
		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
	}, false)
	if err != nil {
		return err
	}
	if _, err := ctx.Write(b); err != nil {
		return err
	}
	return ioutil.PipeBuffer(ctx, <-tf.tunCh, tf.sbuf, tf.dbuf)
}