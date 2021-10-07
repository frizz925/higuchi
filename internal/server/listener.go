package server

import (
	"net"
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/frizz925/higuchi/internal/errors"
	"github.com/frizz925/higuchi/internal/filter"
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
	for l.running.Load() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		logger := l.logger.With(zap.String("src", conn.RemoteAddr().String()))
		logger.Info("Accepted connection")
		l.pool.Dispatch(&filter.Context{
			Conn:   conn,
			Logger: logger,
		}, l.connCallback)
	}
}

func (l *Listener) connCallback(c *filter.Context, err error) {
	logger := c.Logger
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
			logger.Error("Proxy error", zap.Error(err))
		} else {
			res.StatusCode = http.StatusInternalServerError
			logger.Error("Connection error", zap.Error(err))
		}
		b, err := httputil.DumpResponse(res, false)
		if err != nil {
			logger.Error("Failed creating response", zap.Error(err))
		} else if _, err := c.Write(b); err != nil {
			logger.Error("Failed writing response", zap.Error(err))
		}
	}
	if err := c.Close(); err != nil {
		logger.Error("Close connection error", zap.Error(err))
	}
}
