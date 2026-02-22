package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"

	"github.com/mickamy/sampay/internal/lib/validator"
)

type APIConfig struct {
	Port int `env:"PORT" envDefault:"8080" validate:"required"`
}

var (
	apiOnce sync.Once
	api     APIConfig
)

func API() APIConfig {
	apiOnce.Do(func() {
		if err := env.Parse(&api); err != nil {
			panic(err)
		}

		if err := validator.Struct(context.Background(), &api); err != nil {
			panic(fmt.Errorf("invalid api config: %+v", err))
		}
	})

	return api
}
