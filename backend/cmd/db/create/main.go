package main

import (
	"context"
	"fmt"
	"os"

	"mickamy.com/sampay/internal/cli/db"
)

func main() {
	fmt.Println("Creating database...")
	if err := db.Create(context.Background()); err != nil {
		fmt.Println("failed to create database: ", err)
		os.Exit(1)
		return
	}
	fmt.Println("Done.")
}
