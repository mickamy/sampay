package database

import (
	"context"
	"fmt"

	"github.com/mickamy/ormgen/orm"
)

type DB struct {
	orm.Querier

	conn    *orm.DB
	cleanup func()
}

type Option func(*DB)

func WithCleanup(cleanup func()) Option {
	return func(db *DB) {
		db.cleanup = cleanup
	}
}

func New(conn *orm.DB, options ...Option) *DB {
	db := &DB{
		Querier: conn,
		conn:    conn,
	}

	for _, option := range options {
		option(db)
	}

	return db
}

func (db *DB) Transaction(ctx context.Context, fn func(tx *DB) error) error {
	if err := db.conn.Transaction(ctx, func(tx *orm.Tx) error {
		return fn(&DB{Querier: tx, conn: db.conn})
	}); err != nil {
		return err //nolint:wrapcheck // trivial to wrap error here
	}
	return nil
}

func (db *DB) Close() error {
	if db.cleanup != nil {
		db.cleanup()
	}
	if err := db.conn.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	return nil
}

type Writer struct {
	*DB
}

type Reader struct {
	*DB
}

type ReadWriter struct {
	*Reader
	*Writer
}
