package db

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/mickamy/slogger"

	"mickamy.com/sampay/config"
)

func Migrate(ctx context.Context) error {
	cfg := config.Database()

	// migrate
	{
		migrations := filepath.Join(config.Common().PackageRoot, "db", "migrations")
		m, err := migrate.New(fmt.Sprintf("file://%s", migrations), cfg.WriterURL())
		if err != nil {
			return fmt.Errorf("failed to initialize migrations: %w", err)
		}

		defer func(m *migrate.Migrate) {
			if err, _ := m.Close(); err != nil {
				slogger.WarnCtx(ctx, "failed to close database: %v", err)
			}
		}(m)

		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to migrate: %w", err)
		}
	}

	// grant reader to select
	{
		variables := map[string]string{
			"sampay.db_name":         cfg.Name.Escape(),
			"sampay.reader_username": cfg.Reader.Escape(),
		}

		if err := runPSQL(ctx, "03_grant_reader_to_read.sql", cfg.Writer.Escape(), cfg.WriterPass.Escape(), cfg.Name.Escape(), variables); err != nil {
			return fmt.Errorf("failed to grant read user: %w", err)
		}
	}

	return nil
}
