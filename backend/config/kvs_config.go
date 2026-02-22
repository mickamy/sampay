package config

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/caarlos0/env/v11"

	"github.com/mickamy/sampay/internal/lib/validator"
)

type KVSConfig struct {
	Host     string `env:"KVS_HOST" validate:"required"`
	Port     int    `env:"KVS_PORT" validate:"required"`
	Username string `env:"KVS_USERNAME"` // may be empty on development
	Password string `env:"KVS_PASSWORD" validate:"required"`
}

func (c KVSConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
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

		if err := validator.Struct(context.Background(), &kvs); err != nil {
			panic(fmt.Errorf("invalid kvs config: %+v", err))
		}
	})
	return kvs
}
