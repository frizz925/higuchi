package server

import (
	"io"
	"net"
	"sync"

	"github.com/frizz925/higuchi/internal/errors"
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
		c, err := l.Accept()
		if err != nil {
			if err != io.EOF {
				l.logger.Error("failed to accept connection", zap.Error(err))
			}
			continue
		}
		l.logger.Info("accepted connection", zap.String("src", c.RemoteAddr().String()))
		l.pool.Dispatch(c, l.connCallback)
	}
}

func (l *Listener) connCallback(c net.Conn, err error) {
	raddr := c.RemoteAddr().String()
	logger := l.logger.With(zap.String("src", raddr))
	if err != nil {
		logger = logger.With(zap.Error(err))
		if v, ok := err.(*errors.EstablishedError); ok {
			logger.Error("connection error")
		} else {
			logger.Error("proxy error", zap.String("user", v.User.Username()), zap.String("dst", v.Destination.String()))
		}
	}
	if err := c.Close(); err != nil {
		logger.Error("close connection error", zap.Error(err))
	}
}
