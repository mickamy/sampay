package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"

	"github.com/mickamy/sampay/internal/lib/validator"
)

type AuthConfig struct {
	SigningSecret string `env:"JWT_SIGNING_SECRET" validate:"required"`
}

func (c AuthConfig) SigningSecretBytes() []byte {
	return []byte(c.SigningSecret)
}

var (
	authOnce sync.Once
	auth     AuthConfig
)

func Auth() AuthConfig {
	authOnce.Do(func() {
		if err := env.Parse(&auth); err != nil {
			panic(err)
		}

		if err := validator.Struct(context.Background(), &auth); err != nil {
			panic(fmt.Errorf("invalid auth config: %+v", err))
		}
	})

	return auth
}
