package di

import (
	"context"
	"fmt"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/infra/storage/kvs"
)

type Infra struct {
	_        context.Context  `inject:"param"` //nolint:containedctx // required by injector
	_        config.KVSConfig `inject:"provider:config.KVS"`
	DB       *database.DB     `inject:"provider:di.ProvideDB"`
	WriterDB *database.Writer `inject:""`
	ReaderDB *database.Reader `inject:""`
	KVS      *kvs.KVS         `inject:"provider:di.ProvideKVS"`
}

func (i *Infra) Close() error {
	i.KVS.Close()
	if err := i.WriterDB.Close(); err != nil {
		return fmt.Errorf("failed to close writer db: %w", err)
	}
	if err := i.ReaderDB.Close(); err != nil {
		return fmt.Errorf("failed to close reader db: %w", err)
	}

	return nil
}

func ProvideDB(
	ctx context.Context,
	commonCfg config.CommonConfig,
	databaseCfg config.DatabaseConfig,
) (*database.DB, error) {
	writer, err := ProvideWriterDB(ctx, commonCfg, databaseCfg)
	if err != nil {
		return nil, err
	}
	return writer.DB, nil
}

func ProvideWriterDB(
	ctx context.Context,
	commonCfg config.CommonConfig,
	databaseCfg config.DatabaseConfig,
) (*database.Writer, error) {
	db, err := database.Open(ctx, commonCfg, databaseCfg.WriterProvider(), database.RoleWriter)
	if err != nil {
		return nil, fmt.Errorf("failed to open writer db: %w", err)
	}
	return &database.Writer{DB: db}, nil
}

func ProvideReaderDB(
	ctx context.Context,
	commonCfg config.CommonConfig,
	databaseCfg config.DatabaseConfig,
) (*database.Reader, error) {
	db, err := database.Open(ctx, commonCfg, databaseCfg.ReaderProvider(), database.RoleReader)
	if err != nil {
		return nil, fmt.Errorf("failed to open reader db: %w", err)
	}
	return &database.Reader{DB: db}, nil
}

func ProvideKVS(cfg config.KVSConfig) (*kvs.KVS, error) {
	kvStore, err := kvs.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("di: failed to initialize KVS: %w", err)
	}

	return kvStore, nil
}
