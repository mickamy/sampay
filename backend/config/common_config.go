package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
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
