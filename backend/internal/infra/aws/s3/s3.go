package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/mickamy/sampay/config"
)

const defaultPresignExpiry = 15 * time.Minute

type Client interface {
	PresignPutObject(ctx context.Context, bucket, key string) (string, error)
}

type client struct {
	s3        *s3.Client
	presigner *s3.PresignClient
}

func New(ctx context.Context, cfg config.AWSConfig) (Client, error) {
	var opts []func(*aconfig.LoadOptions) error
	if cfg.LocalStackEndpoint != "" {
		opts = append(opts, aconfig.WithBaseEndpoint(cfg.LocalStackEndpoint))
	}

	awsCfg, err := aconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("s3: failed to load aws config: %w", err)
	}

	svc := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.LocalStackEndpoint != "" {
			o.UsePathStyle = true
		}
	})

	return &client{
		s3:        svc,
		presigner: s3.NewPresignClient(svc),
	}, nil
}

func (c *client) PresignPutObject(ctx context.Context, bucket, key string) (string, error) {
	out, err := c.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(defaultPresignExpiry))
	if err != nil {
		return "", fmt.Errorf("s3: failed to presign put object: %w", err)
	}
	return out.URL, nil
}
