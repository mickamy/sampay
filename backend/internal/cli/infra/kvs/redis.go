package kvs

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"

	"mickamy.com/sampay/config"
)

var (
	once     sync.Once
	instance *KVS
)

func Connect(cfg config.KVSConfig) (*KVS, error) {
	var err error
	if instance == nil {
		once.Do(func() {
			var opts *redis.Options
			opts, err = redis.ParseURL(cfg.URL)
			if err != nil {
				err = fmt.Errorf("failed to parse redis url: %s", err)
			}

			instance = redis.NewClient(opts)
		})
	}

	if err != nil {
		return nil, err
	}

	return instance, nil
}
