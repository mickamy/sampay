package database

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mickamy.com/sampay/internal/lib/slices"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

// DB is a wrapper of *gorm.DB
type DB struct{ *gorm.DB }

// Writer is a writer of DB
type Writer struct {
	*DB
}

// Reader is a reader of DB
type Reader struct {
	*DB
}

// WriterTransactional is a transactional for Writer
type WriterTransactional interface {
	// WriterTransaction is a transaction for Writer
	// f is a function that receives a WriterTransactional and returns an error
	WriterTransaction(ctx context.Context, f func(tx WriterTransactional) error) error

	// LockForUpdate is a method that locks the rows for update
	LockForUpdate() WriterTransactional

	// WriterDB is a method that returns a DB of Writer
	WriterDB() *DB
}

func (w Writer) WriterTransaction(ctx context.Context, f func(tx WriterTransactional) error) error {
	return w.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return f(&Writer{&DB{tx}})
	})
}

func (w Writer) LockForUpdate() WriterTransactional {
	withLock := w.Clauses(clause.Locking{Strength: "UPDATE"})
	return &Writer{&DB{withLock}}
}

func (w Writer) WriterDB() *DB {
	return w.DB
}

// ReaderTransactional is a transactional for Reader
type ReaderTransactional interface {
	// ReaderTransaction is a transaction for Reader
	// f is a function that receives a ReaderTransactional and returns an error
	ReaderTransaction(ctx context.Context, f func(tx ReaderTransactional) error) error

	// ReaderDB is a method that returns a DB of Reader
	ReaderDB() *DB
}

func (r Reader) ReaderTransaction(ctx context.Context, f func(tx ReaderTransactional) error) error {
	return r.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return f(&Reader{&DB{tx}})
	})
}

func (r Reader) ReaderDB() *DB {
	return r.DB
}

// ReadWriter is a wrapper of Writer and Reader
type ReadWriter struct {
	writer *Writer
	reader *Reader
}

func NewReadWriter(writer *Writer, reader *Reader) *ReadWriter {
	return &ReadWriter{writer, reader}
}

func (db ReadWriter) WriterTransaction(ctx context.Context, f func(tx WriterTransactional) error) error {
	return db.writer.WriterTransaction(ctx, f)
}

func (db ReadWriter) ReaderTransaction(ctx context.Context, f func(tx ReaderTransactional) error) error {
	return db.reader.ReaderTransaction(ctx, f)
}

func (db ReadWriter) LockForUpdate() WriterTransactional {
	return db.writer.LockForUpdate()
}

func (db ReadWriter) WriterDB() *DB {
	return db.writer.DB
}

func (db ReadWriter) ReaderDB() *DB {
	return db.reader.DB
}

func (db ReadWriter) Writer() *Writer {
	return db.writer
}

func (db ReadWriter) Reader() *Reader {
	return db.reader
}

var (
	_ WriterTransactional = (*ReadWriter)(nil)
	_ ReaderTransactional = (*ReadWriter)(nil)
)

type Scope func(db *DB) *DB

func (s Scope) Gorm() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return s(&DB{db}).DB
	}
}

type Scopes []Scope

func (ss Scopes) Gorm() []func(db *gorm.DB) *gorm.DB {
	return slices.Map(ss, func(s Scope) func(db *gorm.DB) *gorm.DB {
		return s.Gorm()
	})
}
