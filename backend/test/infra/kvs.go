package infra

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
)

func NewKVS(t *testing.T) *kvs.KVS {
	t.Helper()

	if useTestContainers {
		return initRedisContainers(t)
	}

	return initActualRedis(t)
}

func initRedisContainers(t *testing.T) *redis.Client {
	ctx := context.Background()

	ctn, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         uuid.NewString(),
			Image:        "redis:7.0.15-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor: wait.ForAll(
				wait.ForLog("* Ready to accept connections"),
				wait.ForListeningPort("6379/tcp"),
			),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("could not start redis: %s", err)
	}

	host, err := ctn.Host(ctx)
	if err != nil {
		t.Fatalf("could not get host: %s", err)
	}

	port, err := ctn.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("could not get port %s: %s", "6379", err)
	}

	t.Cleanup(func() {
		if err := ctn.Terminate(ctx); err != nil {
			t.Fatalf("could not stop redis: %s", err)
		}
	})

	rds := redis.NewClient(&redis.Options{
		Addr: host + ":" + port.Port(),
	})

	return rds
}

func initActualRedis(t *testing.T) *redis.Client {
	opts, err := redis.ParseURL(config.KVS().URL)
	if err != nil {
		panic(fmt.Errorf("failed to parse redis url: %s", err))
	}

	client := redis.NewClient(opts)

	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Logf("failed to close redis connection: %s", err)
		}
	})

	return client
}
