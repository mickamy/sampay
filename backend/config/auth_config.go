package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type AuthConfig struct {
	SigningSecret string `env:"JWT_SIGNING_SECRET"`
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

		if auth.SigningSecret == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %#v", auth))
		}
	})

	return auth
}
