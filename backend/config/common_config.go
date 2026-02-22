package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"

	"github.com/mickamy/sampay/internal/lib/validator"
)

type Env string

func (e Env) String() string {
	return string(e)
}

func (e Env) ShortName() string {
	switch e {
	case EnvDevelopment:
		return "dev"
	case EnvTest:
		return "test"
	case EnvStaging:
		return "stg"
	case EnvProduction:
		return "prod"
	}

	panic(fmt.Sprintf("unknown environment: %s", e))
}

func (e Env) IsDevelopment() bool {
	return e == EnvDevelopment
}

func (e Env) IsTest() bool {
	return e == EnvTest
}

func (e Env) IsStaging() bool {
	return e == EnvStaging
}

func (e Env) IsProduction() bool {
	return e == EnvProduction
}

func (e Env) ShouldLogToFile() bool {
	return e == EnvStaging || e == EnvProduction
}

const (
	EnvDevelopment Env = "development"
	EnvTest        Env = "test"
	EnvStaging     Env = "staging"
	EnvProduction  Env = "production"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type CommonConfig struct {
	Env        Env      `env:"ENV" validate:"required,oneof=development test staging production"`
	LogLevel   LogLevel `env:"LOG_LEVEL"   envDefault:"debug" validate:"required,oneof=debug info warn error"`
	ModuleRoot string   `env:"MODULE_ROOT" validate:"required"`
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

		if err := validator.Struct(context.Background(), &common); err != nil {
			panic(fmt.Errorf("invalid common config: %+v", err))
		}
	})

	return common
}
