package main

import (
	"fmt"
	"os"

	"mickamy.com/sampay/internal/api"
	"mickamy.com/sampay/internal/di"
)

func main() {
	infras, err := di.InitInfras()
	if err != nil {
		fmt.Println("failed to initialize infras:", err)
		os.Exit(1)
	}
	s := api.NewServer(infras)

	if err := s.ListenAndServe(); err != nil {
		fmt.Println("failed to start server:", err)
		os.Exit(1)
	}
}
