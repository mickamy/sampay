package main

import (
	"context"
	"fmt"
	"os"

	"mickamy.com/sampay/internal/cli/db"
)

func main() {
	fmt.Println("Seeding database...")
	if err := db.Seed(context.Background()); err != nil {
		fmt.Println("failed to seed database: ", err)
		os.Exit(1)
		return
	}
	fmt.Println("Done.")
}
