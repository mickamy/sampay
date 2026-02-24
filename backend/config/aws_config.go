package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"

	"github.com/mickamy/sampay/internal/lib/validator"
)

type AWSConfig struct {
	S3PublicBucket     string `env:"S3_PUBLIC_BUCKET_NAME" validate:"required"`
	S3PrivateBucket    string `env:"S3_PRIVATE_BUCKET_NAME" validate:"required"`
	LocalStackEndpoint string `env:"LOCALSTACK_ENDPOINT"`
	CloudfrontDomain   string `env:"CLOUDFRONT_DOMAIN" validate:"required"`
}

func (c AWSConfig) CloudfrontURL() string {
	return "https://" + c.CloudfrontDomain
}

var (
	awsOnce sync.Once
	awsCfg  AWSConfig
)

func AWS() AWSConfig {
	awsOnce.Do(func() {
		if err := env.Parse(&awsCfg); err != nil {
			panic(err)
		}

		if err := validator.Struct(context.Background(), &awsCfg); err != nil {
			panic(fmt.Errorf("invalid aws config: %+v", err))
		}
	})
	return awsCfg
}
