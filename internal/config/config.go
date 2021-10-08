package config

import (
	"encoding/base64"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func (l Logger) Create() *zap.Logger {
	var cfg zapcore.EncoderConfig
	switch l.Mode {
	case "production":
		cfg = zap.NewProductionEncoderConfig()
	default:
		cfg = zap.NewDevelopmentEncoderConfig()
	}

	var zenc zapcore.Encoder
	switch l.Encoding {
	case "console":
		zenc = zapcore.NewConsoleEncoder(cfg)
	default:
		zenc = zapcore.NewJSONEncoder(cfg)
	}

	lvlEnablerInfo := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if l.Mode == "development" {
			return lvl <= zapcore.InfoLevel
		}
		return lvl == zapcore.InfoLevel
	})
	lvlEnablerErr := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	stdout := zapcore.Lock(os.Stdout)
	stderr := zapcore.Lock(os.Stderr)
	core := zapcore.NewTee(
		zapcore.NewCore(zenc, stdout, lvlEnablerInfo),
		zapcore.NewCore(zenc, stderr, lvlEnablerErr),
	)
	return zap.New(core)
}

func (a Auth) Pepper() ([]byte, error) {
	return base64.StdEncoding.DecodeString(a.pepper)
}
