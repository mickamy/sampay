package database

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DB = gorm.DB

type Transactional interface {
	Transaction(ctx context.Context, f func(tx Transactional) error) error
	ReaderTransaction(ctx context.Context, f func(tx Transactional) error) error
	LockForUpdate() Transactional
	WriterInstance() *DB
	ReaderInstance() *DB
}

type ReadWriter struct {
	Writer *DB
	Reader *DB
}

func (db ReadWriter) Transaction(ctx context.Context, f func(tx Transactional) error) error {
	return db.Writer.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return f(ReadWriter{Writer: tx, Reader: db.Reader})
	})
}

func (db ReadWriter) ReaderTransaction(ctx context.Context, f func(tx Transactional) error) error {
	return db.Reader.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return f(ReadWriter{Writer: db.Writer, Reader: tx})
	})
}

func (db ReadWriter) LockForUpdate() Transactional {
	return ReadWriter{Writer: db.Writer.Clauses(clause.Locking{Strength: "UPDATE"}), Reader: db.Reader}
}

func (db ReadWriter) WriterInstance() *DB {
	return db.Writer
}

func (db ReadWriter) ReaderInstance() *DB {
	return db.Reader
}

var (
	_ Transactional = (*ReadWriter)(nil)
)
