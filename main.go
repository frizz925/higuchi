package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/frizz925/higuchi/internal/dispatcher"
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/pool"
	"github.com/frizz925/higuchi/internal/server"
	"github.com/frizz925/higuchi/internal/worker"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer logger.Sync()

	users := map[string]string{"testuser": "proxytestpass"}
	s := server.New(server.Config{
		Logger: logger,
		Pool: pool.NewPreallocatedPool(func(num int) *worker.Worker {
			return worker.New(num, filter.NewParseFilter(
				filter.NewAuthFilter(users),
				filter.NewTunnelFilter(server.DefaultBufferSize),
				filter.NewForwardFilter(filter.NewDispatchFilter(dispatcher.NewTCPDispatcher(server.DefaultBufferSize))),
			))
		}, 1024),
	})

	if _, err := s.Listen("tcp", "0.0.0.0:8080"); err != nil {
		return err
	}
	logger.Info("Higuchi listening at :8080")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	return s.Close()
}
