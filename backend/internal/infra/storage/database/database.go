package database

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mickamy/ormgen/orm"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/lib/logger"
)

type ConnectionRole string

const (
	RoleWriter ConnectionRole = "writer"
	RoleReader ConnectionRole = "reader"
)

var (
	writerMu sync.Mutex
	readerMu sync.Mutex
	writerDB *DB
	readerDB *DB
)

func Open(
	ctx context.Context,
	cfg config.CommonConfig,
	provider config.DatabaseConfigProvider,
	role ConnectionRole,
) (*DB, error) {
	switch role {
	case RoleWriter:
		if writerDB != nil {
			return writerDB, nil
		}

		writerMu.Lock()
		defer writerMu.Unlock()

		if writerDB != nil {
			return writerDB, nil
		}

		instance, err := open(ctx, cfg, provider)
		if err != nil {
			return nil, err
		}

		writerDB = instance
		return writerDB, nil

	case RoleReader:
		if readerDB != nil {
			return readerDB, nil
		}

		readerMu.Lock()
		defer readerMu.Unlock()

		if readerDB != nil {
			return readerDB, nil
		}

		instance, err := open(ctx, cfg, provider)
		if err != nil {
			return nil, err
		}

		readerDB = instance
		return readerDB, nil
	}

	return nil, fmt.Errorf("unexpected connection type: %s", role)
}

func open(
	ctx context.Context, cfg config.CommonConfig, provider config.DatabaseConfigProvider,
) (*DB, error) {
	var drvName = "pgx"

	sqlDB, err := sql.Open(drvName, provider.URL())
	if err != nil {
		return nil, fmt.Errorf("database: failed to open database: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database: failed to ping database: %w", err)
	}

	conn := orm.New(sqlDB, orm.PostgreSQL)
	if cfg.LogLevel.ShouldLogORM() {
		conn = conn.Debug(logger.ORM{})
	}

	return New(conn), nil
}
