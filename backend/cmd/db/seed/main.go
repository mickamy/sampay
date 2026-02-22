package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/misc/seed"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	envs := os.Getenv("SEED_ENV")
	if envs == "" {
		envs = config.Common().Env.String()
	}
	env := config.Env(envs)

	cfg := config.Common()
	cfg.Env = env

	fmt.Printf("Seeding data for environment %s...\n", env)

	ctx := context.Background()
	db, err := database.Open(ctx, cfg, config.Database().WriterProvider(), database.RoleWriter)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() { _ = db.Close() }()

	if err := seed.Do(ctx, &database.Writer{DB: db}, env); err != nil {
		return fmt.Errorf("failed to seed data: %w", err)
	}

	fmt.Println("Done.")
	return nil
}
