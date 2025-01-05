package infra

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/db/seed"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/lib/slices"
)

var (
	seedOnce = sync.Once{}
)

type CleanUp = func()

type WriterDSN string
type ReaderDSN string

type DSN struct {
	Writer WriterDSN
	Reader ReaderDSN
}

func NewDB() (DSN, CleanUp) {
	if useTestContainers {
		return initPostgresContainers(config.Database())
	}

	return initActualDB(config.Database())
}

func initPostgresContainers(cfg config.DatabaseConfig) (DSN, CleanUp) {
	ctx := context.Background()

	packageRoot := config.Common().PackageRoot
	mountFiles := slices.Map([]string{
		"00_users.sql",
		"01_database.sql",
		"02_migrate.sh",
		"03_grant_select_to_reader.sql",
	}, func(file string) testcontainers.ContainerFile {
		return testcontainers.ContainerFile{
			HostFilePath:      filepath.Join(packageRoot, "db", file),
			ContainerFilePath: "/docker-entrypoint-initdb.d/" + file,
			FileMode:          0o644,
		}
	})
	ctn, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         uuid.NewString(),
			Image:        "postgres:16.4-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "postgres",
			},
			HostConfigModifier: func(hostConfig *container.HostConfig) {
				migrations := mount.Mount{
					Type:   mount.TypeBind,
					Source: filepath.Join(packageRoot, "db", "migrations"),
					Target: "/docker-entrypoint-initdb.d/migrations",
				}
				hostConfig.Mounts = append(hostConfig.Mounts, migrations)
			},
			Files: mountFiles,
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
				wait.ForListeningPort("5432/tcp"),
			).WithDeadline(60 * time.Second),
		},
		Started: true,
		Reuse:   reuseContainer,
	})
	if err != nil {
		log.Fatalf("cloud not start postgres: %s", err)
	}

	host, err := ctn.Host(ctx)
	if err != nil {
		log.Fatalf("cloud not get host: %s", err)
	}
	port, err := ctn.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("cloud not get port %s: %s", "5432", err)
	}

	writerDSN := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s TimeZone=%s sslmode=disable",
		host,
		port.Int(),
		cfg.Writer,
		cfg.Name,
		cfg.WriterPass,
		cfg.TimeZone,
	)
	writerDB, err := gorm.Open(postgres.New(postgres.Config{DSN: writerDSN}), &gorm.Config{})
	if err != nil {
		log.Fatalf("cloud not connect to writer database: %s", err)
	}
	readerDSN := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s TimeZone=%s sslmode=disable",
		host,
		port.Int(),
		cfg.Reader,
		cfg.Name,
		cfg.ReaderPass,
		cfg.TimeZone,
	)

	readerDB, err := gorm.Open(postgres.New(postgres.Config{DSN: readerDSN}), &gorm.Config{})
	if err != nil {
		log.Fatalf("cloud not connect to reader database: %s", err)
	}

	if err := seed.Do(ctx, &database.Writer{DB: &database.DB{DB: writerDB}}, "test"); err != nil {
		log.Fatalf("failed to seed: %s", err)
	}

	return DSN{WriterDSN(writerDSN), ReaderDSN(readerDSN)}, func() {
		for _, db := range []*gorm.DB{writerDB, readerDB} {
			sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("cloud not get DB connection: %s", err)
			}
			if err := sqlDB.Close(); err != nil {
				log.Fatalf("cloud not close DB connection: %s", err)
			}
			if err := ctn.Terminate(ctx); err != nil {
				log.Fatalf("cloud not terminate container: %s", err)
			}
		}
	}
}

func initActualDB(cfg config.DatabaseConfig) (DSN, CleanUp) {
	writerDSN := cfg.WriterDSN()
	writer, err := gorm.Open(postgres.New(postgres.Config{DSN: writerDSN}))
	if err != nil {
		log.Fatalf("failed to connect to writer database: %s", err)
	}

	readerDSN := cfg.ReaderDSN()
	reader, err := gorm.Open(postgres.New(postgres.Config{DSN: readerDSN}))
	if err != nil {
		log.Fatalf("failed to connect to reader database: %s", err)
	}

	seedOnce.Do(func() {
		ctx := context.Background()
		if err := seed.Do(ctx, &database.Writer{DB: &database.DB{DB: writer}}, "test"); err != nil {
			log.Fatalf("failed to seed: %s", err)
		}
	})

	return DSN{WriterDSN(writerDSN), ReaderDSN(readerDSN)}, func() {
		for _, db := range []*gorm.DB{writer, reader} {
			sqlDB, err := db.DB()
			if err != nil {
				log.Fatalf("cloud not get DB connection: %s", err)
			}
			err = sqlDB.Close()
			if err != nil {
				log.Fatalf("cloud not close DB connection: %s", err)
			}
		}
	}
}

func OpenTXDB(t *testing.T, dsn string) *database.DB {
	t.Helper()

	driverName := "txdb_" + t.Name()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DriverName: driverName,
		DSN:        dsn,
	}))
	if err != nil {
		t.Fatalf("failed to execute gorm.Open: %s", err)
	}

	t.Cleanup(func() {
		sqlDB, err := gormDB.DB()
		if err != nil {
			t.Logf("failed to get database instance: %s", err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Logf("failed to close database: %s", err)
		}
	})

	return &database.DB{DB: gormDB}
}
