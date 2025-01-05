package infra

import (
	"context"
	"fmt"
	"log"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type KVSAddr string

type CleanUpKVS = func()

func NewKVS() (KVSAddr, CleanUpKVS) {
	if useTestContainers {
		return initRedisContainers()
	}

	return initMiniRedis()
}

func initRedisContainers() (KVSAddr, CleanUpKVS) {
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
		log.Fatalf("could not start redis: %s", err)
	}

	host, err := ctn.Host(ctx)
	if err != nil {
		log.Fatalf("could not get host: %s", err)
	}

	port, err := ctn.MappedPort(ctx, "6379")
	if err != nil {
		log.Fatalf("could not get port %s: %s", "6379", err)
	}

	addr := fmt.Sprintf("%s:%s", host, port.Port())

	return KVSAddr(addr), func() {
		if err := ctn.Terminate(ctx); err != nil {
			log.Fatalf("could not stop redis: %s", err)
		}
	}
}

func initMiniRedis() (KVSAddr, CleanUpKVS) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("could not start miniredis: %s", err)
	}

	return KVSAddr(mr.Addr()), func() {
		mr.Close()
	}
}
