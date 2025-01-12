package s3

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"mickamy.com/sampay/config"
	sdkConfig "mickamy.com/sampay/internal/lib/aws/config"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type Client interface {
	GeneratePresignedURL(ctx context.Context, bucket, key string, secs int64) (string, error)

	PutObject(ctx context.Context, bucket, key string, body io.Reader) error

	GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)
}

type client struct {
	client *s3.Client
}

func New(cfg config.AWSConfig) Client {
	sdkCfg := sdkConfig.Load(context.Background())
	c := s3.NewFromConfig(sdkCfg, func(o *s3.Options) {
		if cfg.LocalstackEndpoint != "" {
			o.BaseEndpoint = aws.String(cfg.LocalstackEndpoint)
			o.UsePathStyle = true
		}
	})

	return &client{c}
}

func (c *client) GeneratePresignedURL(ctx context.Context, bucket, key string, secs int64) (string, error) {
	client := s3.NewPresignClient(c.client)

	req, err := client.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(time.Duration(secs)*time.Second))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned url: %w", err)
	}

	return req.URL, nil
}

func (c *client) PutObject(ctx context.Context, bucket, key string, body io.Reader) error {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}

	return nil
}

func (c *client) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	resp, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return resp.Body, nil
}
