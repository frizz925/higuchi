package config

import (
	"encoding/base64"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Config struct {
	Server struct {
		Listeners []string
	}
	Pool struct {
		WorkerCount int
	}
	Worker struct {
		BufferSize int
	}
	Logger  Logger
	Filters struct {
		Auth        Auth
		Certbot     Certbot
		Forwarded   Forwarded
		Healthcheck Healthcheck
	}
}

type Logger struct {
	Mode              string
	Encoding          string
	DisableCaller     bool
	DisableStackTrace bool
}

type Auth struct {
	Enabled       bool
	PasswordsFile string
	pepper        string `mapstructure:"Pepper"`
}

type Certbot struct {
	Enabled       bool
	Hostname      string
	Webroot       string
	ChallengePath string
}

type Forwarded struct {
	Enabled bool
}

type Healthcheck struct {
	Enabled bool
	Method  string
	Path    string
}

func ReadConfig() (cfg Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&cfg); err != nil {
		return
	}
	return
}

func (l Logger) Create() (*zap.Logger, error) {
	var zc zap.Config
	switch l.Mode {
	case "production":
		zc = zap.NewProductionConfig()
	default:
		zc = zap.NewDevelopmentConfig()
	}
	zc.Encoding = l.Encoding
	zc.DisableCaller = l.DisableCaller
	zc.DisableStacktrace = l.DisableStackTrace
	return zc.Build()
}

func (a Auth) Pepper() ([]byte, error) {
	return base64.StdEncoding.DecodeString(a.pepper)
}
