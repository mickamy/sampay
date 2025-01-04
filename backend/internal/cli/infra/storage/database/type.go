package database

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Database = gorm.DB

type Transactional interface {
	Transaction(ctx context.Context, f func(tx Transactional) error) error
	ReaderTransaction(ctx context.Context, f func(tx Transactional) error) error
	LockForUpdate() Transactional
	Writer() *Database
	Reader() *Database
}

type TransactionalDatabase struct {
	writer *Database
	reader *Database
}

func (db TransactionalDatabase) Transaction(ctx context.Context, f func(tx Transactional) error) error {
	return db.writer.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return f(TransactionalDatabase{writer: tx, reader: db.reader})
	})
}

func (db TransactionalDatabase) ReaderTransaction(ctx context.Context, f func(tx Transactional) error) error {
	return db.reader.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return f(TransactionalDatabase{writer: db.writer, reader: tx})
	})
}

func (db TransactionalDatabase) LockForUpdate() Transactional {
	return TransactionalDatabase{writer: db.writer.Clauses(clause.Locking{Strength: "UPDATE"}), reader: db.reader}
}

func (db TransactionalDatabase) Writer() *Database {
	return db.writer
}

func (db TransactionalDatabase) Reader() *Database {
	return db.reader
}

var (
	_ Transactional = (*TransactionalDatabase)(nil)
)
