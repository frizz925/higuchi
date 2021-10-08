package filter

import (
	"bufio"
	"net"
	"net/http"
	"sync"

	"github.com/frizz925/higuchi/internal/httputil"
	"github.com/frizz925/higuchi/internal/ioutil"
	"github.com/frizz925/higuchi/internal/netutil"
)

type TunnelFilter struct {
	ioPool sync.Pool
}

type tunnelIO struct {
	sbuf, dbuf []byte
	bw         *bufio.Writer
}

// Tunnel filter is a filter to intercept tunneling proxy request.
// The filter has its own chain when intercepting a tunneling connection.
// Otherwise, it would just use the provided chain instead.
func NewTunnelFilter(bufsize int) *TunnelFilter {
	tf := &TunnelFilter{}
	tf.ioPool.New = func() interface{} {
		return &tunnelIO{
			sbuf: make([]byte, bufsize),
			dbuf: make([]byte, bufsize),
			bw:   bufio.NewWriterSize(nil, bufsize),
		}
	}
	return tf
}

func (tf *TunnelFilter) Do(ctx *Context, req *http.Request, next Next) error {
	if req.Method != http.MethodConnect {
		return next()
	}
	// Create the connection to target host
	conn, err := net.Dial("tcp", httputil.ParseRequestAddress(req))
	if err != nil {
		return err
	}
	ctx.Logger.Info("Established connection for tunneling")
	// Remove the prefixed buffer if we're tunneling
	if v, ok := ctx.Conn.(*netutil.PrefixedConn); ok {
		ctx.Conn = v.Conn
	}

	tio := tf.ioPool.Get().(*tunnelIO)
	defer tf.ioPool.Put(tio)
	tio.bw.Reset(ctx)
	httputil.WriteResponseHeader(&http.Response{
		StatusCode: http.StatusOK,
		Status:     "200 Connection established",
		Proto:      req.Proto,
		ProtoMajor: req.ProtoMajor,
		ProtoMinor: req.ProtoMinor,
	}, tio.bw)
	if err := tio.bw.Flush(); err != nil {
		return err
	}
	return ioutil.PipeBuffer(ctx.Conn, conn, tio.sbuf, tio.dbuf)
}
