package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/caarlos0/env/v11"
)

type AuthConfig struct {
	EmailVerificationExpiresIn int    `env:"EMAIL_VERIFICATION_EXPIRES_IN" envDefault:"86400"` // seconds
	SigningSecret              string `env:"JWT_SIGNING_SECRET"`
}

func (c AuthConfig) EmailVerificationExpiresInDuration() time.Duration {
	return time.Duration(c.EmailVerificationExpiresIn) * time.Second
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
