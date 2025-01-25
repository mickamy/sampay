package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

func initBySSM(ctx context.Context, region string, env Env) error {
	cfg, err := config.LoadDefaultConfig(ctx, func(options *config.LoadOptions) error {
		options.Region = region
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to load SDK config: %w", err)
	}

	client := ssm.NewFromConfig(cfg)

	ssmPathPrefix := fmt.Sprintf("/sampay/app/%s/", env.ShortName())
	paginator := ssm.NewGetParametersByPathPaginator(client, &ssm.GetParametersByPathInput{
		Path:           aws.String(ssmPathPrefix),
		WithDecryption: aws.Bool(true),
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to get parameters: %w", err)
		}
		if err := setEnvVar(output.Parameters, ssmPathPrefix); err != nil {
			return fmt.Errorf("failed to set environment variables: %w", err)
		}
	}

	return nil
}

func setEnvVar(parameters []types.Parameter, pathPrefix string) error {
	for _, parameter := range parameters {
		envVarName := strings.TrimPrefix(*parameter.Name, pathPrefix)

		if err := os.Setenv(envVarName, *parameter.Value); err != nil {
			return err
		}
	}

	return nil
}
