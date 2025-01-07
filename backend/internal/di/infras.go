package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
)

type Infras struct {
	*database.DB
	*database.ReadWriter
	*database.Writer
	*database.Reader
	*kvs.KVS
}

func NewInfras(readWriter *database.ReadWriter, kvs *kvs.KVS) Infras {
	return Infras{
		DB:         readWriter.WriterDB(),
		ReadWriter: readWriter,
		Writer:     readWriter.Writer(),
		Reader:     readWriter.Reader(),
		KVS:        kvs,
	}
}

func provideDB(cfg config.DatabaseConfig) (*database.DB, error) {
	writer, err := provideWriter(cfg)
	if err != nil {
		return nil, err
	}
	return &database.DB{DB: writer.DB.DB}, nil
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

	return &database.Writer{DB: writer}, nil
}

func provideReader(cfg config.DatabaseConfig) (*database.Reader, error) {
	reader, err := database.Connect(cfg.ReaderProvider(), database.ReaderConnection)
	if err != nil {
		return nil, err
	}

	return &database.Reader{DB: reader}, nil
}

func provideKVS(cfg config.KVSConfig) (*kvs.KVS, error) {
	return kvs.Connect(cfg)
}

//lint:ignore U1000 used by wire
var infras = wire.NewSet(
	provideDB,
	provideReadWriter,
	provideWriter,
	provideReader,
	provideKVS,
)
