package server

import (
	"bufio"
	"net"
	"net/http"
	"sync"

	"github.com/frizz925/higuchi/internal/errors"
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/httputil"
	"github.com/frizz925/higuchi/internal/pool"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

type Listener struct {
	net.Listener

	pool           pool.Pool
	logger         *zap.Logger
	running        atomic.Bool
	wg             sync.WaitGroup
	removeListener func(*Listener)
}

func (l *Listener) Close() error {
	if err := l.stop(); err != nil {
		return err
	}
	l.removeListener(l)
	return nil
}

func (l *Listener) start() {
	l.running.Store(true)
	l.wg.Add(1)
	go l.runRoutine()
}

func (l *Listener) stop() error {
	l.running.Store(false)
	if err := l.Listener.Close(); err != nil {
		return err
	}
	l.wg.Wait()
	return nil
}

func (l *Listener) runRoutine() {
	defer l.wg.Done()
	laddr := l.Addr().String()
	for l.running.Load() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		raddr := conn.RemoteAddr().String()
		ctx := filter.NewContext(conn, l.logger, filter.LogFields{
			Proto:    "http",
			Listener: laddr,
			Source:   raddr,
		})
		ctx.Logger.Info("Accepted connection")
		l.pool.Dispatch(ctx, l.connCallback)
	}
}

func (l *Listener) connCallback(c *filter.Context, err error) {
	if err != nil {
		res := &http.Response{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		if v, ok := err.(*errors.HTTPError); ok {
			res.Proto = v.Request.Proto
			res.ProtoMajor = v.Request.ProtoMajor
			res.ProtoMinor = v.Request.ProtoMinor
			res.StatusCode = v.StatusCode
			res.Header = v.Header
			c.Logger.Error("Proxy error", zap.Error(err))
		} else {
			res.StatusCode = http.StatusInternalServerError
			c.Logger.Error("Connection error", zap.Error(err))
		}
		bw := bufio.NewWriter(c)
		httputil.WriteResponseHeader(res, bw)
		if err := bw.Flush(); err != nil {
			c.Logger.Error("Failed writing response", zap.Error(err))
		}
	}
	if err := c.Close(); err != nil {
		c.Logger.Error("Close connection error", zap.Error(err))
	}
}
