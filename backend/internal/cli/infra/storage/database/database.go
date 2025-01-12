package database

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/lib/logger"
)

type ConnectionType string

const (
	WriterConnection ConnectionType = "default"
	ReaderConnection ConnectionType = "default_reader"
)

var (
	writerOnce sync.Once
	readerOnce sync.Once

	writerDB *DB
	readerDB *DB
)

func Connect(provider config.DatabaseConfigProvider, connType ConnectionType) (*DB, error) {
	var err error
	switch connType {
	case WriterConnection:
		if writerDB == nil {
			writerOnce.Do(func() {
				writerDB, err = initializeDB(provider)
			})
		}
		return writerDB, err
	case ReaderConnection:
		if readerDB == nil {
			readerOnce.Do(func() {
				readerDB, err = initializeDB(provider)
			})
		}
		return readerDB, err
	default:
		return nil, fmt.Errorf("unexpected connection type: %s", connType)
	}
}
func initializeDB(provider config.DatabaseConfigProvider) (*DB, error) {
	db, err := gorm.Open(
		postgres.New(postgres.Config{DSN: provider.DSN()}),
		&gorm.Config{Logger: logger.Gorm.LogMode(logLevel())},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &DB{db}, nil
}

func logLevel() gormLogger.LogLevel {
	switch config.Common().LogLevel {
	case "debug", "info":
		return gormLogger.Info
	case "warn":
		return gormLogger.Warn
	case "error":
		return gormLogger.Error
	}
	return gormLogger.Info
}
