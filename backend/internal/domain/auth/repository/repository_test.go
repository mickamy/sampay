package repository_test

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
	kvsAddr     infra.KVSAddr
)

func TestMain(m *testing.M) {
	var wg sync.WaitGroup
	wg.Add(2)

	databaseDSNCh := make(chan infra.DatabaseDSN)
	cleanUpDBCh := make(chan func())

	kvsAddrCh := make(chan infra.KVSAddr)
	cleanUpKVSCh := make(chan func())

	go func() {
		defer wg.Done()
		dsn, c := infra.NewDB()
		databaseDSNCh <- dsn
		cleanUpDBCh <- c
	}()

	go func() {
		defer wg.Done()
		addr, c := infra.NewKVS()
		kvsAddrCh <- addr
		cleanUpKVSCh <- c
	}()

	databaseDSN = <-databaseDSNCh
	cleanUpDatabase := <-cleanUpDBCh

	kvsAddr = <-kvsAddrCh
	cleanUpKVS := <-cleanUpKVSCh

	wg.Wait()

	defer cleanUpDatabase()
	defer cleanUpKVS()

	os.Exit(m.Run())
}

func newReadWriter(t *testing.T) *database.ReadWriter {
	t.Helper()
	txdb := infra.OpenTXDB(t, string(databaseDSN.Writer))
	return database.NewReadWriter(&database.Writer{DB: txdb}, &database.Reader{DB: txdb})
}

func newKVS(t *testing.T) *kvs.KVS {
	t.Helper()
	return redis.NewClient(&redis.Options{
		Addr: string(kvsAddr),
	})
}
