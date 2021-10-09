package server

import (
	"bufio"
	"net"
	"net/http"
	"sync"

	intErrors "github.com/frizz925/higuchi/internal/errors"
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/httputil"
	"github.com/frizz925/higuchi/internal/pool"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

const DefaultBufferSize = 1024

type Server struct {
	net.Listener

	pool         pool.Pool
	logger       *zap.Logger
	running      atomic.Bool
	wg           sync.WaitGroup
	removeServer func(*Server)
}

type Config struct {
	Pool   pool.Pool
	Logger *zap.Logger
}

func Serve(l net.Listener, cfg Config) *Server {
	srv := &Server{
		Listener: l,
		pool:     cfg.Pool,
		logger:   cfg.Logger,
	}
	srv.start()
	return srv
}

func (l *Server) Close() error {
	if err := l.stop(); err != nil {
		return err
	}
	if l.removeServer != nil {
		l.removeServer(l)
	}
	return nil
}

func (l *Server) start() {
	l.running.Store(true)
	l.wg.Add(1)
	go l.runRoutine()
}

func (l *Server) stop() error {
	l.running.Store(false)
	if err := l.Listener.Close(); err != nil {
		return err
	}
	l.wg.Wait()
	return nil
}

func (l *Server) runRoutine() {
	defer l.wg.Done()
	laddr := l.Addr().String()
	for l.running.Load() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		raddr := conn.RemoteAddr().String()
		ctx := filter.NewContext(conn, l.logger, filter.LogFields{
			Proto:  "http",
			Server: laddr,
			Source: raddr,
		})
		ctx.Logger.Info("Accepted connection")
		l.pool.Dispatch(ctx, l.connCallback)
	}
}

func (l *Server) connCallback(c *filter.Context, err error) {
	if err != nil {
		res := &http.Response{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
		}
		if v, ok := err.(*intErrors.HTTPError); ok {
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
