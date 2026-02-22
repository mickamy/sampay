package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"

	"github.com/mickamy/sampay/internal/lib/validator"
)

type OAuthConfig struct {
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID" validate:"required"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET" validate:"required"`
	LINEChannelID      string `env:"LINE_CHANNEL_ID" validate:"required"`
	LINEChannelSecret  string `env:"LINE_CHANNEL_SECRET" validate:"required"`
	RedirectURL        string `env:"OAUTH_REDIRECT_URL" validate:"required"`
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

		if err := validator.Struct(context.Background(), &oauth); err != nil {
			panic(fmt.Errorf("invalid oauth config: %+v", err))
		}
	})

	return oauth
}
