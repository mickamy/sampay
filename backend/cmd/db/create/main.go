package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/lib/pq"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/lib/either"
)

func main() {
	cfg := config.Database()

	db := dbmate.New(either.Must(url.Parse(cfg.AdminURL())))
	if err := db.Create(); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "42P04" {
				log.Println("database already exists")
				os.Exit(0)
			}
		}
		log.Fatalf("failed to create database: %v", err)
	}

	fmt.Println("Done.")
}
