package db

import (
	"context"
	"fmt"

	"mickamy.com/sampay/config"
)

func Create(ctx context.Context) error {
	db := config.Database()
	common := config.Common()

	// Create users
	if common.Env == config.Development {
		variables := map[string]string{
			"sampay.writer_username": db.Writer,
			"sampay.writer_password": db.WriterPass,
			"sampay.reader_username": db.Reader,
			"sampay.reader_password": db.ReaderPass,
		}
		if err := runPSQL("00_users.sql", variables); err != nil {
			return fmt.Errorf("failed to create users: %w", err)
		}
	}

	// Create database
	{
		variables := map[string]string{
			"sampay.db_name":         db.Name,
			"sampay.writer_username": db.Writer,
		}
		if err := runPSQL("01_database.sql", variables); err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}
