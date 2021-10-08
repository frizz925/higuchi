package server

import (
	"errors"
	"net"
	"os"
	"sync"

	"github.com/frizz925/higuchi/internal/pool"
	"go.uber.org/zap"
)

const DefaultBufferSize = 1024

var (
	ErrServerAlreadyRunning = errors.New("server already running")
	ErrServerNotRunning     = errors.New("server not running")
)

type Server struct {
	pool      pool.Pool
	logger    *zap.Logger
	listeners map[*Listener]struct{}
	mu        sync.Mutex
}

type Config struct {
	Pool   pool.Pool
	Logger *zap.Logger
}

func New(cfg Config) *Server {
	return &Server{
		pool:      cfg.Pool,
		logger:    cfg.Logger,
		listeners: make(map[*Listener]struct{}),
	}
}

func (s *Server) Listen(network string, address string) (*Listener, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	isUnix := network == "unix"
	if isUnix {
		if err := s.maybeRemoveUnixSocket(address); err != nil {
			return nil, err
		}
	}
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	if isUnix {
		if err := s.fixUnixSocketPermissions(address); err != nil {
			l.Close()
			return nil, err
		}
	}

	ls := &Listener{
		Listener:       l,
		pool:           s.pool,
		logger:         s.logger,
		removeListener: s.removeListener,
	}
	ls.start()
	s.listeners[ls] = struct{}{}
	return ls, nil
}

func (s *Server) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for l := range s.listeners {
		if err := l.stop(); err != nil {
			return err
		}
		delete(s.listeners, l)
	}
	return nil
}

func (s *Server) maybeRemoveUnixSocket(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(name)
}

func (s *Server) fixUnixSocketPermissions(name string) error {
	return os.Chmod(name, 0666)
}

func (s *Server) removeListener(l *Listener) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.listeners, l)
}
