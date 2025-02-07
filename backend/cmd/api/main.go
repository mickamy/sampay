package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mickamy/slogger"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/api"
	"mickamy.com/sampay/internal/di"
)

func init() {
	cfg := config.Common()
	slogger.Init(slogger.Config{
		Level:          cfg.SLoggerLevel(),
		TrimPathPrefix: cfg.PackageRoot,
		ContextFieldsExtractor: func(ctx context.Context) []any {
			return []any{}
		},
	})
}

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
