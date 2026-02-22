package config

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

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
	if onDocker {
		loadSecrets()
		return
	}

	switch env {
	case EnvStaging, EnvProduction:
		region, found := os.LookupEnv("AWS_REGION")
		if !found {
			panic("AWS_REGION environment variable is not set")
		}

		panic(fmt.Sprintf("TODO: fetch environment values from SecretsManager or SSM: region=[%s]", region))
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

var secrets = []string{
	"google_client_id",
	"google_client_secret",
	"line_channel_id",
	"line_channel_secret",
}

func loadSecrets() {
	for _, name := range secrets {
		p := path.Join("/run/secrets", name)
		b, err := os.ReadFile(p) //nolint:gosec // path is constructed from hardcoded base path and secret names
		if err != nil {
			continue
		}
		envKey := strings.ToUpper(name)
		if err := os.Setenv(envKey, strings.TrimSpace(string(b))); err != nil {
			panic(fmt.Errorf("failed to set env var %s: %w", envKey, err))
		}
	}
}
