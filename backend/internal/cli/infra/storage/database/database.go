package database

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mickamy.com/sampay/config"
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
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: provider.DSN()}))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return &DB{db}, nil
}
