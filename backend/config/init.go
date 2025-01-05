package config

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	var env Env
	if strings.HasSuffix(os.Args[0], ".test") {
		env = Test
	} else {
		envStr, _ := os.LookupEnv("ENV")
		env = Env(envStr)
		if env == "" {
			env = Development
		}
	}

	if err := os.Setenv("ENV", string(env)); err != nil {
		panic(fmt.Errorf("failed to set env var: %s", err))
	}

	var onDocker bool
	if docker, ok := os.LookupEnv("DOCKER"); ok {
		onDocker = docker == "1"
	}
	if onDocker {
		return
	}

	switch env {
	case "staging", "production":
		region, found := os.LookupEnv("AWS_REGION")
		if !found {
			panic("AWS_REGION environment variable is not set")
		}
		if err := initBySSM(context.Background(), region, env); err != nil {
			panic(fmt.Errorf("failed to init env from SSM: %w", err))
		}
	case "development", "test":
		packageRoot, ok := os.LookupEnv("PACKAGE_ROOT")
		if !ok {
			panic("PACKAGE_ROOT environment variable is not set")
		}
		envFile := path.Join(packageRoot, ".env")
		if err := godotenv.Load(envFile); err != nil {
			panic(fmt.Errorf("error loading .env file: %w", err))
		}
	default:
		panic(fmt.Errorf("unknown ENV variable %s", env))
	}
}
