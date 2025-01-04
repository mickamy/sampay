package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/mickamy/slogger"

	"mickamy.com/sampay/config"
)

func Drop(ctx context.Context) error {
	if config.Common().Env == config.Production {
		return fmt.Errorf("cannot drop database in production")
	}

	cfg := config.Database()
	db, err := sql.Open("postgres", cfg.AdminDSN())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slogger.Warn("failed to close DB connection", "err", err)
		}
	}(db)
	if _, err := db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s (FORCE)", cfg.Name)); err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	return nil
}
