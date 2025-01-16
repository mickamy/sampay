package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type WEBConfig struct {
	BaseURL string `env:"FRONTEND_BASE_URL"`
}

var (
	webOnce sync.Once
	web     WEBConfig
)

func WEB() WEBConfig {
	webOnce.Do(func() {
		if err := env.Parse(&web); err != nil {
			panic(err)
		}

		if web.BaseURL == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %#v", web))
		}
	})

	return web
}
