package seed

import (
	"context"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

type seed func(ctx context.Context, writer *database.Writer, env config.Env) error

func Do(ctx context.Context, writer *database.Writer, env config.Env) error {
	for _, fn := range []seed{} {
		if err := fn(ctx, writer, env); err != nil {
			return err
		}
	}

	return nil
}
