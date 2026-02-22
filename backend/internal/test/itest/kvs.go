package itest

import (
	"strconv"
	"testing"

	"github.com/alicebob/miniredis/v2"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/infra/storage/kvs"
	"github.com/mickamy/sampay/internal/lib/either"
)

type cleanupKVS func()

func NewKVS(t *testing.T) *kvs.KVS {
	t.Helper()

	cfg, cu := initMiniRedis(t)
	t.Cleanup(cu)

	return newKVS(t, cfg)
}

func newKVS(t *testing.T, cfg config.KVSConfig) *kvs.KVS {
	t.Helper()

	c, err := kvs.New(cfg, kvs.WithDisableCache())
	if err != nil {
		t.Fatalf("could not create kvs client: %s", err)
	}
	return c
}

func initMiniRedis(t *testing.T) (config.KVSConfig, cleanupKVS) {
	t.Helper()

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("could not start miniredis: %s", err)
	}

	return config.KVSConfig{
		Host: mr.Host(),
		Port: either.Must(strconv.Atoi(mr.Port())),
	}, mr.Close
}
