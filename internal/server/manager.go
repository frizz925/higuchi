package server

import (
	"net"
	"os"
	"sync"

	"github.com/frizz925/higuchi/internal/pool"
	"go.uber.org/zap"
)

type Manager struct {
	pool    pool.Pool
	logger  *zap.Logger
	servers map[*Server]struct{}
	mu      sync.Mutex
}

type ManagerConfig struct {
	Pool   pool.Pool
	Logger *zap.Logger
}

func NewManager(cfg ManagerConfig) *Manager {
	return &Manager{
		pool:    cfg.Pool,
		logger:  cfg.Logger,
		servers: make(map[*Server]struct{}),
	}
}

func (m *Manager) ListenAndServe(network string, address string) (*Server, error) {
	isUnix := network == "unix"
	if isUnix {
		if err := m.maybeRemoveUnixSocket(address); err != nil {
			return nil, err
		}
	}
	l, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	if isUnix {
		if err := m.fixUnixSocketPermissions(address); err != nil {
			l.Close()
			return nil, err
		}
	}
	return m.Serve(l), nil
}

func (m *Manager) Serve(l net.Listener) *Server {
	m.mu.Lock()
	defer m.mu.Unlock()
	srv := Serve(l, Config{
		Pool:   m.pool,
		Logger: m.logger,
	})
	srv.removeServer = m.removeServer
	m.servers[srv] = struct{}{}
	return srv
}

func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for s := range m.servers {
		if err := s.stop(); err != nil {
			return err
		}
		delete(m.servers, s)
	}
	return nil
}

func (m *Manager) maybeRemoveUnixSocket(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(name)
}

func (m *Manager) fixUnixSocketPermissions(name string) error {
	return os.Chmod(name, 0666)
}

func (m *Manager) removeServer(s *Server) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.servers, s)
}
