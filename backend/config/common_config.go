package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
	"github.com/mickamy/slogger"
)

type Env string

const (
	Development Env = "development"
	Test        Env = "test"
	Staging     Env = "staging"
	Production  Env = "production"
)

type CommonConfig struct {
	Env         Env    `env:"ENV"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	PackageRoot string `env:"PACKAGE_ROOT"`
}

func (c CommonConfig) SLoggerLevel() slogger.Level {
	switch c.LogLevel {
	case "debug":
		return slogger.LevelDebug
	case "info":
		return slogger.LevelInfo
	case "warn":
		return slogger.LevelWarn
	case "error":
		return slogger.LevelError
	default:
		return slogger.LevelInfo
	}
}

var (
	commonOnce sync.Once
	common     CommonConfig
)

func Common() CommonConfig {
	commonOnce.Do(func() {
		if err := env.Parse(&common); err != nil {
			panic(err)
		}

		if common.Env == "" || common.LogLevel == "" || common.PackageRoot == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %#v", common))
		}
	})

	return common
}
