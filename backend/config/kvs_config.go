package config

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/caarlos0/env/v11"
)

type KVSConfig struct {
	Host     string `env:"KVS_HOST" envDefault:"localhost"`
	Port     int    `env:"KVS_PORT" envDefault:"6379"`
	Password string `env:"KVS_PASSWORD"`
}

func (c KVSConfig) URL() string {
	return fmt.Sprintf("redis://:%s@%s:%d", url.QueryEscape(c.Password), c.Host, c.Port)
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

		if kvs.Host == "" || kvs.Port == 0 {
			panic(fmt.Errorf("some of required environment variables are missing: %+v", kvs))
		}
	})
	return kvs
}
