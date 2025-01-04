package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type KVSConfig struct {
	URL string `env:"REDIS_URL"`
}

var (
	kvsOnce sync.Once
	kvs     KVSConfig
)

func KVS() KVSConfig {
	kvsOnce.Do(func() {
		if err := env.Parse(&kvs); err != nil {
			panic(err)
		}

		if kvs.URL == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %+v", kvs))
		}
	})
	return kvs
}
