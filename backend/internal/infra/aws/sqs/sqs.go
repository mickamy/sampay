package sqs

import (
	"context"
	"fmt"

	aconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"github.com/mickamy/sampay/config"
)

func New(ctx context.Context, cfg config.AWSConfig) (*sqs.Client, error) {
	var opts []func(*aconfig.LoadOptions) error
	if cfg.LocalStackEndpoint != "" {
		opts = append(opts, aconfig.WithBaseEndpoint(cfg.LocalStackEndpoint))
	}

	awsCfg, err := aconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("sqs: failed to load aws config: %w", err)
	}

	return sqs.NewFromConfig(awsCfg), nil
}
