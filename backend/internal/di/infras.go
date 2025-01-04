package di

import (
	"errors"

	"github.com/google/wire"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/kvs"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
)

type Infras struct {
	*database.ReadWriter
	*kvs.KVS
}

func provideDB(cfg config.DatabaseConfig) (*database.ReadWriter, error) {
	var errs []error

	writer, err := database.Connect(cfg, database.WriterConnection)
	if err != nil {
		errs = append(errs, err)
	}

	reader, err := database.Connect(cfg, database.ReaderConnection)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return &database.ReadWriter{Writer: writer, Reader: reader}, nil
}

func provideKVS(cfg config.KVSConfig) (*kvs.KVS, error) {
	return kvs.Connect(cfg)
}

//lint:ignore U1000 used by wire
var infras = wire.NewSet(
	provideDB,
	provideKVS,
)
