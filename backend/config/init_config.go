package config

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"slices"

	"github.com/joho/godotenv"
)

func init() {
	var env Env
	if isTest() {
		env = EnvTest
	} else {
		envStr := os.Getenv("ENV")
		env = Env(envStr)
		if env == "" {
			env = EnvDevelopment
		}
	}

	if err := os.Setenv("ENV", string(env)); err != nil {
		panic(fmt.Errorf("failed to set env var: %w", err))
	}

	var onDocker bool
	if docker, ok := os.LookupEnv("DOCKER"); ok {
		onDocker = docker == "1"
	}
	if onDocker && slices.Contains([]Env{EnvDevelopment, EnvTest}, env) {
		return
	}

	switch env {
	case EnvStaging, EnvProduction:
		region, found := os.LookupEnv("AWS_REGION")
		if !found {
			panic("AWS_REGION environment variable is not set")
		}

		if err := initBySecretsManager(context.Background(), region, env.SecretID()); err != nil {
			panic(fmt.Errorf("failed to init secrets manager: %w", err))
		}
	case EnvDevelopment, EnvTest:
		moduleRoot, ok := os.LookupEnv("MODULE_ROOT")
		if !ok {
			panic("MODULE_ROOT environment variable is not set")
		}
		envFile := path.Join(moduleRoot, ".env")
		if err := godotenv.Load(envFile); err != nil {
			panic(fmt.Errorf("error loading .env file: %w", err))
		}
	default:
		panic(fmt.Errorf("unknown ENV variable %s", env))
	}
}

func isTest() bool {
	return flag.Lookup("test.v") != nil ||
		flag.Lookup("test.run") != nil ||
		flag.Lookup("test.count") != nil
}
