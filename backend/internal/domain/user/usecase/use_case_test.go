package usecase_test

import (
	"os"
	"sync"
	"testing"

	"github.com/redis/go-redis/v9"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
	"mickamy.com/sampay/test/infra"
)

var (
	databaseDSN infra.DatabaseDSN
)

func TestMain(m *testing.M) {
	var wg sync.WaitGroup
	wg.Add(1)

	databaseDSNCh := make(chan infra.DatabaseDSN)
	cleanUpDBCh := make(chan func())

	go func() {
		defer wg.Done()
		dsn, c := infra.NewDB()
		databaseDSNCh <- dsn
		cleanUpDBCh <- c
	}()

	databaseDSN = <-databaseDSNCh
	cleanUpDatabase := <-cleanUpDBCh

	wg.Wait()

	defer cleanUpDatabase()

	os.Exit(m.Run())
}

func newReadWriter(t *testing.T) *database.ReadWriter {
	t.Helper()
	txdb := infra.OpenTXDB(t, string(databaseDSN.Writer))
	return database.NewReadWriter(&database.Writer{DB: txdb}, &database.Reader{DB: txdb})
}

func newKVS(t *testing.T) *kvs.KVS {
	t.Helper()

	addr, c := infra.NewKVS()
	t.Cleanup(c)

	return redis.NewClient(&redis.Options{
		Addr: string(addr),
	})
}
