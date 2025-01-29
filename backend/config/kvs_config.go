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
	Password string `env:"KVS_PASSWORD" envDefault:""`
}

func (c KVSConfig) URL() string {
	scheme := "redis"

	// Encode password if present
	var kvsAuth string
	if c.Password != "" {
		kvsAuth = url.QueryEscape(c.Password) + "@"
	}

	return fmt.Sprintf("%s://%s%s:%d", scheme, kvsAuth, c.Host, c.Port)
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
