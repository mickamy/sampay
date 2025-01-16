package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type EmailConfig struct {
	From string `env:"EMAIL_FROM"`
}

var (
	emailOnce sync.Once
	email     EmailConfig
)

func Email() EmailConfig {
	emailOnce.Do(func() {
		if err := env.Parse(&email); err != nil {
			panic(err)
		}

		if email.From == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %#v", email))
		}
	})

	return email
}
