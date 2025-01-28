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
			"sampay.writer_username": db.Writer.Escape(),
			"sampay.writer_password": db.WriterPass.Escape(),
			"sampay.reader_username": db.Reader.Escape(),
			"sampay.reader_password": db.ReaderPass.Escape(),
		}
		if err := runPSQL("00_users.sql", variables); err != nil {
			return fmt.Errorf("failed to create users: %w", err)
		}
	}

	// Create database
	{
		variables := map[string]string{
			"sampay.db_name":         db.Name.Escape(),
			"sampay.writer_username": db.Writer.Escape(),
		}
		if err := runPSQL("01_database.sql", variables); err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}
