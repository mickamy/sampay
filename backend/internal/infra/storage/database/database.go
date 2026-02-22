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
	writerOnce    sync.Once
	readerOnce    sync.Once
	writerDB      *DB
	readerDB      *DB
	writerOpenErr error //nolint:errname // not a sentinel error; used for sync.Once result caching
	readerOpenErr error //nolint:errname // not a sentinel error; used for sync.Once result caching
)

func Open(
	ctx context.Context,
	cfg config.CommonConfig,
	provider config.DatabaseConfigProvider,
	role ConnectionRole,
) (*DB, error) {
	switch role {
	case RoleWriter:
		writerOnce.Do(func() {
			writerDB, writerOpenErr = open(ctx, cfg, provider)
		})
		return writerDB, writerOpenErr

	case RoleReader:
		readerOnce.Do(func() {
			readerDB, readerOpenErr = open(ctx, cfg, provider)
		})
		return readerDB, readerOpenErr
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
