package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type AWSConfig struct {
	CloudFrontDomain string `env:"CLOUDFRONT_DOMAIN"`
}

var (
	awsOnce sync.Once
	_aws    AWSConfig
)

func AWS() AWSConfig {
	awsOnce.Do(func() {
		if err := env.Parse(&_aws); err != nil {
			panic(err)
		}

		if _aws.CloudFrontDomain == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %#v", _aws))
		}
	})

	return _aws
}
