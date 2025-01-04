package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
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
		m, err := migrate.New(fmt.Sprintf("file://%s", migrations), cfg.URL())
		if err != nil {
			return fmt.Errorf("failed to initialize migrations: %w", err)
		}

		defer func(m *migrate.Migrate) {
			err, _ := m.Close()
			if err != nil {
				slogger.WarnCtx(ctx, "failed to close database: %v", err)
			}
		}(m)

		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to migrate: %w", err)
		}
	}

	// grant reader to select
	{
		cmd := exec.Command("psql", "-U", cfg.AdminUser, "-h", cfg.Host, "-d", "postgres", "-f", path.Join(config.Common().PackageRoot, "db", "03_grant_select_to_reader.sql"))
		cmd.Env = append(os.Environ(), "PGPASSWORD="+cfg.AdminPassword)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to grant read user: %w\nOutput: %s", err, string(output))
		}
	}

	return nil
}
