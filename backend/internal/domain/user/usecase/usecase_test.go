package usecase_test

import (
	"os"
	"testing"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/test/itest"
)

var databaseDSN itest.DatabaseDSN

func TestMain(m *testing.M) {
	dsn, cleanup := itest.NewDB()
	databaseDSN = dsn

	code := m.Run()
	cleanup()
	os.Exit(code)
}

func newReadWriter(t *testing.T) *database.ReadWriter {
	t.Helper()
	txdb := itest.OpenTXDB(t, string(databaseDSN.Writer))
	return &database.ReadWriter{Reader: &database.Reader{DB: txdb}, Writer: &database.Writer{DB: txdb}}
}

func newInfra(t *testing.T) *di.Infra {
	t.Helper()
	readWriter := newReadWriter(t)
	return &di.Infra{
		DB:       readWriter.Writer.DB,
		WriterDB: readWriter.Writer,
		ReaderDB: readWriter.Reader,
		KVS:      itest.NewKVS(t),
	}
}
