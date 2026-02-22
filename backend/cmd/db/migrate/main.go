package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/lib/either"
)

func main() {
	cfg := config.Database()

	db := dbmate.New(either.Must(url.Parse(cfg.AdminURL())))
	if err := db.Migrate(); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	fmt.Println("Done.")
}
