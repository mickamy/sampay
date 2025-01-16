package ses

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"

	"mickamy.com/sampay/config"
	sdkConfig "mickamy.com/sampay/internal/lib/aws/config"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type Client interface {
	Send(ctx context.Context, input *sesv2.SendEmailInput) (*sesv2.SendEmailOutput, error)
}

type client struct {
	client *sesv2.Client
}

func New(cfg config.AWSConfig) Client {
	sdkCfg := sdkConfig.Load(context.Background())
	c := sesv2.NewFromConfig(sdkCfg, func(o *sesv2.Options) {
		if cfg.SESEndpoint != "" {
			o.BaseEndpoint = aws.String(cfg.SESEndpoint)
		}
	})

	return &client{c}
}

func (c *client) Send(ctx context.Context, input *sesv2.SendEmailInput) (*sesv2.SendEmailOutput, error) {
	return c.client.SendEmail(ctx, input)
}
