package config

import (
	"fmt"
	"sync"

	"github.com/caarlos0/env/v11"
)

type AWSConfig struct {
	AccessKeyID                 string `env:"AWS_ACCESS_KEY_ID"`
	AccessKeySecret             string `env:"AWS_ACCESS_KEY_SECRET"`
	CloudFrontDomain            string `env:"CLOUDFRONT_DOMAIN"`
	LocalstackEndpoint          string `env:"LOCALSTACK_ENDPOINT"`
	Region                      string `env:"AWS_REGION"`
	S3PublicBucket              string `env:"S3_PUBLIC_BUCKET_NAME"`
	SESEndpoint                 string `env:"SES_ENDPOINT"`
	SQSWorkerURL                string `env:"SQS_WORKER_URL"`
	SQSWorkerDeadLetterQueueURL string `env:"SQS_WORKER_DLQ_URL"`
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

		if _aws.CloudFrontDomain == "" || _aws.Region == "" || _aws.S3PublicBucket == "" {
			panic(fmt.Errorf("some of required environment variables are missing: %#v", _aws))
		}
	})

	return _aws
}
