package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/frizz925/higuchi/internal/pool"
	"github.com/frizz925/higuchi/internal/server"
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

	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		return err
	}
	defer listener.Close()

	pool := pool.NewFixedPool(pool.FixedPoolConfig{
		Logger: logger,
	})
	server := server.New(server.Config{
		Logger:   logger,
		Listener: listener,
		Pool:     pool,
	})
	if err := server.Start(); err != nil {
		return err
	}
	defer server.Stop()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	sig := <-ch
	logger.Info("received signal", zap.String("signal", sig.String()))

	return nil
}
