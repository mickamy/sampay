package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/mickamy/go-sqs-worker/consumer"
	sjob "github.com/mickamy/go-sqs-worker/job"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/infra/aws/sqs"
	"github.com/mickamy/sampay/internal/job"
	"github.com/mickamy/sampay/internal/lib/logger"
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
	defer func() {
		if closeErr := infra.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "failed to close infrastructure: %s\n", closeErr)
		}
	}()

	awsCfg := config.AWS()
	kvsCfg := config.KVS()

	sqsClient, err := sqs.New(ctx, awsCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize SQS client: %w", err)
	}

	redisURL := fmt.Sprintf("redis://:%s@%s", kvsCfg.Password, kvsCfg.Address())
	jobs := job.NewJobs(infra)

	c, err := consumer.New(consumer.Config{
		WorkerQueueURL:     awsCfg.SQSWorkerURL,
		DeadLetterQueueURL: awsCfg.SQSWorkerDLQURL,
		RedisURL:           redisURL,
	}, sqsClient, nil, func(jobType string) (sjob.Job, error) {
		return job.Get(jobType, jobs)
	})
	if err != nil {
		return fmt.Errorf("failed to initialize consumer: %w", err)
	}

	workersCount := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			logger.Info(ctx, fmt.Sprintf("worker %d starting", workerID))
			c.Do(ctx)
			logger.Info(ctx, fmt.Sprintf("worker %d finished", workerID))
		}(i)
	}

	wg.Wait()

	return nil
}
