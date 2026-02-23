package repository_test

import (
	"os"
	"sync"
	"testing"

	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/test/itest"
)

var (
	databaseDSN itest.DatabaseDSN
)

func TestMain(m *testing.M) {
	var wg sync.WaitGroup
	wg.Add(1)

	databaseDSNCh := make(chan itest.DatabaseDSN)
	cleanUpDBCh := make(chan func())

	go func() {
		defer wg.Done()
		dsn, c := itest.NewDB()
		databaseDSNCh <- dsn
		cleanUpDBCh <- c
	}()

	databaseDSN = <-databaseDSNCh
	cleanUpDatabase := <-cleanUpDBCh

	wg.Wait()

	code := m.Run()
	cleanUpDatabase()
	os.Exit(code)
}

func newReadWriter(t *testing.T) *database.ReadWriter {
	t.Helper()
	txdb := itest.OpenTXDB(t, string(databaseDSN.Writer))
	return &database.ReadWriter{Reader: &database.Reader{DB: txdb}, Writer: &database.Writer{DB: txdb}}
}
