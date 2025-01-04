package main

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/db"
)

func main() {
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
	if err := db.Drop(context.Background()); err != nil {
		panic(err)
	}
	fmt.Println("Done.")
}
