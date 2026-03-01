package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"

	"github.com/mickamy/sampay/internal/lib/validator"
)

type LineMessagingConfig struct {
	ChannelSecret string `env:"LINE_MESSAGING_CHANNEL_SECRET" validate:"required"`
	ChannelToken  string `env:"LINE_MESSAGING_CHANNEL_TOKEN"  validate:"required"`
}

var (
	lineMessagingOnce sync.Once
	lineMessaging     LineMessagingConfig
)

func LineMessaging() LineMessagingConfig {
	lineMessagingOnce.Do(func() {
		if err := env.Parse(&lineMessaging); err != nil {
			panic(err)
		}

		if err := validator.Struct(context.Background(), &lineMessaging); err != nil {
			panic(fmt.Errorf("invalid line messaging config: %+v", err))
		}
	})

	return lineMessaging
}
