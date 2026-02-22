package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mickamy/sampay/internal/api"
	"github.com/mickamy/sampay/internal/di"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	infra, err := di.NewInfra(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize infrastructure: %w", err)
	}

	srv := api.NewServer(infra)
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServe()
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//nolint:contextcheck // parent ctx is already canceled
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to shutdown gracefully: %s\n", err)
	}

	if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server error: %w", err)
	}

	if err := infra.Close(); err != nil {
		return fmt.Errorf("failed to close infrastructure: %w", err)
	}

	return nil
}
