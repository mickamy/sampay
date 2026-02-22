package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lib/pq"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/lib/either"
)

func main() {
	if config.Common().Env == "production" {
		fmt.Println("Dropping the database in production is not allowed.")
		os.Exit(1)
	}

	fmt.Println("Dropping database...")
	fmt.Println("This operation is irreversible.")
	fmt.Println("Are you sure you want to continue? (yes/no)")
	var answer string
	if _, err := fmt.Scanln(&answer); err != nil {
		panic(err)
	}
	if answer != "yes" {
		fmt.Println("Aborted.")
		return
	}

	cfg := config.Database()
	ctx := context.Background()

	db := dbmate.New(either.Must(url.Parse(cfg.AdminURL())))
	if err := db.Drop(); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "55006" {
				if err := forceDrop(ctx); err != nil {
					log.Fatalf("failed to force drop: %s", err)
				}
				fmt.Println("Done.")
				return
			}
		}
		log.Fatalf("failed to drop database: %v", err)
	}

	fmt.Println("Done.")
}

func forceDrop(ctx context.Context) error {
	cfg := config.Database()
	db, err := sql.Open("pgx", cfg.AdminMaintenanceURL())
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() { _ = db.Close() }()

	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s WITH (FORCE)", cfg.Name.Escape())
	_, err = db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop database: %w", err)
	}

	return nil
}
