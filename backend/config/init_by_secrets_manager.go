package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func initBySecretsManager(ctx context.Context, region, id string) error {
	cfg, err := config.LoadDefaultConfig(ctx, func(options *config.LoadOptions) error {
		options.Region = region
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to load SDK config: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	out, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(id),
	})
	if err != nil {
		return fmt.Errorf("failed to get secret value: %w", err)
	}
	if out.SecretString == nil {
		return fmt.Errorf("secret string is nil: id=%s", id)
	}

	secretString := *out.SecretString
	var secrets map[string]string
	if err := json.Unmarshal([]byte(secretString), &secrets); err != nil {
		return fmt.Errorf("failed to unmarshal secret value: %w", err)
	}

	for k, v := range secrets {
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("failed to set environment variable: %w", err)
		}
	}

	return nil
}
