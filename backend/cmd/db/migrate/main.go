package main

import (
	"context"
	"fmt"
	"os"

	"mickamy.com/sampay/internal/cli/db"
)

func main() {
	fmt.Println("Migrating database...")
	if err := db.Migrate(context.Background()); err != nil {
		fmt.Println("failed to migrate database: ", err)
		os.Exit(1)
		return
	}
	fmt.Println("Done.")
}
