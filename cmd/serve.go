package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/frizz925/higuchi/internal/auth"
	"github.com/frizz925/higuchi/internal/config"
	"github.com/frizz925/higuchi/internal/crypto/hasher"
	"github.com/frizz925/higuchi/internal/dispatcher"
	"github.com/frizz925/higuchi/internal/filter"
	"github.com/frizz925/higuchi/internal/pool"
	"github.com/frizz925/higuchi/internal/server"
	"github.com/frizz925/higuchi/internal/worker"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const unixAddressPrefix = "unix:"

var unixAddressPrefixLength = len(unixAddressPrefix)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web proxy server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.ReadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error while reading config: %v", err)
			os.Exit(1)
		}
		logger := cfg.Logger.Create()
		defer logger.Sync()
		if err := runServe(cfg, logger); err != nil {
			logger.Fatal("Server error", zap.Error(err))
		}
	},
}

func runServe(cfg config.Config, logger *zap.Logger) error {
	hfs := []filter.HTTPFilter{}
	if cfg.Filters.Forwarded.Enabled {
		logger.Info("Forwarded header filter enabled")
		hfs = append(hfs, filter.DefaultForwardedFilter)
	}
	if cfg.Filters.Certbot.Enabled {
		logger.Info("Certbot filter enabled")
		hfs = append(hfs, filter.NewCertbotFilter(filter.CertbotConfig{
			Hostname:      cfg.Filters.Certbot.Hostname,
			Webroot:       cfg.Filters.Certbot.Webroot,
			ChallengePath: cfg.Filters.Certbot.ChallengePath,
		}))
	}
	if cfg.Filters.Auth.Enabled {
		logger.Info("Auth filter enabled")
		af, err := createAuthFilter(cfg.Filters.Auth)
		if err != nil {
			return err
		}
		hfs = append(hfs, af)
	}
	df := filter.NewDispatchFilter(dispatcher.NewTCPDispatcher(cfg.Worker.BufferSize))
	hfs = append(
		hfs,
		filter.NewTunnelFilter(cfg.Worker.BufferSize),
		filter.NewForwardFilter(df),
	)

	filters := []filter.Filter{}
	if cfg.Filters.Healthcheck.Enabled {
		logger.Info("Healthcheck filter enabled")
		filters = append(filters, filter.NewHealthCheckFilter(
			cfg.Filters.Healthcheck.Method,
			cfg.Filters.Healthcheck.Path,
		))
	}
	filters = append(filters, filter.NewParseFilter(cfg.Worker.BufferSize, hfs...))

	s := server.New(server.Config{
		Logger: logger,
		Pool: pool.NewDynamicPool(func(num int) *worker.Worker {
			return worker.New(num, filters...)
		}),
	})

	for _, addr := range cfg.Server.Listeners {
		network := "tcp"
		if strings.HasPrefix(addr, unixAddressPrefix) {
			network = "unix"
			addr = addr[unixAddressPrefixLength:]
		}
		if _, err := s.Listen(network, addr); err != nil {
			return err
		}
		logger.Info(fmt.Sprintf("Higuchi listening at %s", addr))
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	sig := <-sigCh
	logger.Info("Received stop signal", zap.String("signal", sig.String()))

	return s.Close()
}

func createAuthFilter(authCfg config.Auth) (*filter.AuthFilter, error) {
	pepper, err := authCfg.Pepper()
	if err != nil {
		return nil, fmt.Errorf("error while decoding pepper: %v", err)
	}
	fa := auth.NewFileAuth(hasher.NewMD5Hasher(pepper))
	au, err := fa.ReadPasswordsFile(authCfg.PasswordsFile)
	if err != nil {
		return nil, fmt.Errorf("error while reading passwords file: %v", err)
	}
	users := make(map[string]interface{})
	for user, ad := range au {
		users[user] = ad
	}
	return filter.NewAuthFilter(users, func(password string, i interface{}) bool {
		switch v := i.(type) {
		case string:
			return password == v
		case hasher.PasswordDigest:
			return v.Compare(password) == 0
		}
		return false
	}), nil
}
