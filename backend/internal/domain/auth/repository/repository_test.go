package repository_test

import (
	"os"
	"sync"
	"testing"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/test/infra"
)

var (
	dsn infra.DSN
)

func TestMain(m *testing.M) {
	var wg sync.WaitGroup
	wg.Add(1)

	dsnCh := make(chan infra.DSN)
	cleanUpCh := make(chan func())

	go func() {
		defer wg.Done()
		dsn, c := infra.NewDB()
		dsnCh <- dsn
		cleanUpCh <- c
	}()

	dsn = <-dsnCh
	cleanUp := <-cleanUpCh

	wg.Wait()

	defer cleanUp()

	os.Exit(m.Run())
}

func WriterDB(t *testing.T) *database.DB {
	txdb := infra.OpenTXDB(t, string(dsn.Writer))
	return txdb
}

func ReaderDB(t *testing.T) *database.DB {
	txdb := infra.OpenTXDB(t, string(dsn.Reader))
	return txdb
}
