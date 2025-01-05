package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/kvs"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
)

type Infras struct {
	*database.ReadWriter
	*database.Writer
	*database.Reader
	*kvs.KVS
}

func provideReadWriter(cfg config.DatabaseConfig) (*database.ReadWriter, error) {
	var errs []error
	writer, err := provideWriter(cfg)
	if err != nil {
		errs = append(errs, err)
	}

	reader, err := provideReader(cfg)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errs[0]
	}

	return database.NewReadWriter(writer, reader), nil
}

func provideWriter(cfg config.DatabaseConfig) (*database.Writer, error) {
	writer, err := database.Connect(cfg.WriterProvider(), database.WriterConnection)
	if err != nil {
		return nil, err
	}

	return (*database.Writer)(writer), nil
}

func provideReader(cfg config.DatabaseConfig) (*database.Reader, error) {
	reader, err := database.Connect(cfg.ReaderProvider(), database.ReaderConnection)
	if err != nil {
		return nil, err
	}

	return (*database.Reader)(reader), nil
}

func provideKVS(cfg config.KVSConfig) (*kvs.KVS, error) {
	return kvs.Connect(cfg)
}

//lint:ignore U1000 used by wire
var infras = wire.NewSet(
	provideReadWriter,
	provideWriter,
	provideReader,
	provideKVS,
)
