package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/mickamy/slogger"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/lib/logger"
)

func init() {
	cfg := config.Common()
	writer, err := logger.FileWriter()
	if err != nil {
		fmt.Println("failed to create log file writer:", err)
		os.Exit(1)
	}

	slogger.Init(slogger.Config{
		Level:          cfg.SLoggerLevel(),
		Outputs:        []io.Writer{os.Stdout, writer},
		TrimPathPrefix: cfg.PackageRoot,
		ContextFieldsExtractor: func(ctx context.Context) []any {
			return []any{}
		},
	})
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	workersCount := runtime.GOMAXPROCS(0)

	var wg sync.WaitGroup

	consumer := di.InitConsumers().Consumer
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			slogger.InfoCtx(ctx, fmt.Sprintf("worker %d starting", workerID))
			consumer.Do(ctx)
			slogger.InfoCtx(ctx, fmt.Sprintf("worker %d finished", workerID))
		}(i)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	slogger.InfoCtx(ctx, "shutdown signal received, cancelling context")
	cancel()

	wg.Wait()
	slogger.InfoCtx(ctx, "all workers have finished, shutting down")
}
