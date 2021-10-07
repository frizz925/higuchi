package config

import (
	"encoding/base64"

	"github.com/spf13/viper"
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
	Logger struct {
		Mode string
	}
	Filters struct {
		Auth    Auth
		Certbot Certbot
	}
}

type Auth struct {
	Enabled       bool
	PasswordsFile string
	pepper        string `mapstructure:"Pepper"`
}

type Certbot struct {
	Enabled  bool
	Hostname string
	Webroot  string
}

func (a Auth) Pepper() ([]byte, error) {
	return base64.StdEncoding.DecodeString(a.pepper)
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
