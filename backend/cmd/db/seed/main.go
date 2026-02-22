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
		fmt.Printf("failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = db.Close() }()

	if err := seed.Do(ctx, &database.Writer{DB: db}, env); err != nil {
		fmt.Printf("failed to seed data: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Done.")
}
