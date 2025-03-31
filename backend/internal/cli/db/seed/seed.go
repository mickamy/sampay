package seed

import (
	"context"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/infra/storage/database"
)

func Do(ctx context.Context, writer *database.Writer, env config.Env) error {
	funcs := []func(context.Context, *database.Writer, config.Env) error{
		seedUsageCategory,
		seedUserLinkProvider,
		seedUser,
		seedNotification,
	}

	for _, f := range funcs {
		if err := f(ctx, writer, env); err != nil {
			return err
		}
	}

	return nil
}
