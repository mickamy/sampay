package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"mickamy.com/sampay/config"
)

func Create(ctx context.Context) error {
	cfg := config.Database()

	// Create users
	if config.Common().Env == config.Development {
		variables := map[string]string{
			"sampay.writer_username": cfg.Writer.Escape(),
			"sampay.writer_password": cfg.WriterPass.Escape(),
			"sampay.reader_username": cfg.Reader.Escape(),
			"sampay.reader_password": cfg.ReaderPass.Escape(),
		}
		if err := runPSQL(ctx, "00_create_users.sql", cfg.AdminUser.Escape(), cfg.AdminPass.Escape(), "postgres", variables); err != nil {
			return fmt.Errorf("failed to create users: %w", err)
		}
	}

	// Create database
	{
		adminDB, err := sql.Open("postgres", cfg.AdminDSN())
		if err != nil {
			return fmt.Errorf("failed to open database connection: %w", err)
		}
		defer func() {
			if err := adminDB.Close(); err != nil {
				fmt.Printf("failed to close DB connection: %v\n", err)
			}
		}()

		var exists bool
		query := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)"
		if err := adminDB.QueryRowContext(ctx, query, cfg.Name).Scan(&exists); err != nil {
			return fmt.Errorf("failed to check database existence: %w", err)
		}

		if !exists {
			createDBQuery := fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(cfg.Name.Escape()))
			if _, err := adminDB.ExecContext(ctx, createDBQuery); err != nil {
				return fmt.Errorf("failed to create database: %w", err)
			}
			fmt.Println("Database created successfully")
		} else {
			fmt.Println("Database already exists. Skipping creation.")
		}

	}

	// Grant writer to write
	{
		variables := map[string]string{
			"sampay.writer_username": cfg.Writer.Escape(),
		}
		if err := runPSQL(ctx, "01_grant_writer_to_writer.sql", cfg.AdminUser.Escape(), cfg.AdminPass.Escape(), cfg.Name.Escape(), variables); err != nil {
			return fmt.Errorf("failed to grant write user: %w", err)
		}
	}

	return nil
}
