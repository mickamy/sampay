package config

import (
	"context"
	"fmt"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/aws"

	"mickamy.com/sampay/config"
)

func Load(ctx context.Context) aws.Config {
	sdkCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to load SDK config: %w", err))
	}

	cfg := config.AWS()
	if cfg.AccessKeyID != "" && cfg.AccessKeySecret != "" {
		sdkCfg.Credentials = credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.AccessKeySecret, "")
	}

	sdkCfg.Region = config.AWS().Region

	return sdkCfg
}
