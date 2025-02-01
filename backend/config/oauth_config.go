package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type OAuthConfig struct {
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
}

var (
	oauthOnce sync.Once
	oauth     OAuthConfig
)

func OAuth() OAuthConfig {
	oauthOnce.Do(func() {
		if err := env.Parse(&oauth); err != nil {
			panic(err)
		}

		if oauth.GoogleClientID == "" || oauth.GoogleClientSecret == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %#v", oauth))
		}
	})

	return oauth
}
