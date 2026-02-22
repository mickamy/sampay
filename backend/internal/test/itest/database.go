package itest

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mickamy/ormgen/orm"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/misc/seed"
)

type CleanupDB = func()

type WriterDSN string
type ReaderDSN string

type DatabaseDSN struct {
	Writer WriterDSN
	Reader ReaderDSN
}

func NewDB() (DatabaseDSN, CleanupDB) {
	cfg := config.Database()
	return initActualDB(cfg)
}

func initActualDB(cfg config.DatabaseConfig) (DatabaseDSN, CleanupDB) {
	writerDSN := cfg.WriterDSN()
	writerSQL, err := sql.Open("pgx", writerDSN)
	if err != nil {
		log.Fatalf("failed to connect to writer database: %v", err)
	}
	writer := orm.New(writerSQL, orm.PostgreSQL)

	readerDSN := cfg.ReaderDSN()
	readerSQL, err := sql.Open("pgx", readerDSN)
	if err != nil {
		log.Fatalf("failed to connect to reader database: %v", err)
	}
	reader := orm.New(readerSQL, orm.PostgreSQL)

	ctx := context.Background()
	if err := seed.Do(ctx, &database.Writer{DB: database.New(writer)}, config.EnvTest); err != nil {
		log.Fatalf("failed to seed: %s", err)
	}

	return DatabaseDSN{WriterDSN(writerDSN), ReaderDSN(readerDSN)}, func() {
		for _, db := range []*orm.DB{writer, reader} {
			err := db.Close()
			if err != nil {
				log.Fatalf("cloud not close DB connection: %s", err)
			}
		}
	}
}

func OpenTXDB(t *testing.T, dsn string) *database.DB {
	t.Helper()

	baseDriver := stdlib.GetDefaultDriver()
	sql.Register(t.Name(), baseDriver)

	drvName := uuid.NewString()
	txdb.Register(drvName, t.Name(), dsn)

	sqlDB, err := sql.Open(drvName, dsn)
	if err != nil {
		t.Fatalf("failed to open database: %s", err)
	}

	db := orm.New(sqlDB, orm.PostgreSQL)

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close DB connection: %s", err)
		}
	})

	return database.New(db.Debug(logger.ORM{}))
}
