package db

import (
	"context"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/db/seed"
	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/lib/either"
)

func Seed(ctx context.Context) error {
	return seed.Do(ctx, either.Must(di.InitInfras()).Writer, config.Common().Env)
}
