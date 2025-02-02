package di

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/wire"
	"github.com/mickamy/go-sqs-worker/consumer"
	sqsWorkerJob "github.com/mickamy/go-sqs-worker/job"
	"github.com/mickamy/go-sqs-worker/message"
	"github.com/mickamy/go-sqs-worker/producer"
	"github.com/mickamy/slogger"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/job"
	sdkConfig "mickamy.com/sampay/internal/lib/aws/config"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/ptr"
)

//lint:ignore U1000 used by wire
var jobSet = wire.NewSet(
	job.NewSendEmail,
)

type Producers struct {
	*producer.Producer
}

func provideProducerConfig(aws config.AWSConfig, kvs config.KVSConfig) producer.Config {
	return producer.Config{
		WorkerQueueURL: aws.SQSWorkerURL,
		RedisURL:       kvs.URL(),
		BeforeProduceFunc: func(ctx context.Context, msg message.Message) {
			slogger.InfoCtx(ctx, "producing message", "msg", msg)
		},
		AfterProduceFunc: func(ctx context.Context, msg message.Message) {
			slogger.InfoCtx(ctx, "produced message", "msg", msg)
		},
	}
}

func provideSQSClient(aws config.AWSConfig) *sqs.Client {
	sdkCfg := sdkConfig.Load(context.Background())
	c := sqs.NewFromConfig(sdkCfg, func(o *sqs.Options) {
		if aws.LocalstackEndpoint != "" {
			o.BaseEndpoint = ptr.Of(aws.LocalstackEndpoint)
		}
	})
	return c
}

func provideProducer(cfg producer.Config, client *sqs.Client) *producer.Producer {
	return either.Must(producer.New(cfg, client))
}

//lint:ignore U1000 used by wire
var producerSet = wire.NewSet(
	provideProducerConfig,
	provideSQSClient,
	provideProducer,
)

type Consumers struct {
	*consumer.Consumer
}

func provideConsumerConfig(aws config.AWSConfig, kvs config.KVSConfig) consumer.Config {
	return consumer.Config{
		WorkerQueueURL:     aws.SQSWorkerURL,
		DeadLetterQueueURL: aws.SQSWorkerDeadLetterQueueURL,
		RedisURL:           kvs.URL(),
		BeforeProcessFunc: func(ctx context.Context, msg message.Message) error {
			slogger.InfoCtx(ctx, "processing message", "msg", msg)
			return nil
		},
		AfterProcessFunc: func(ctx context.Context, output consumer.Output) error {
			if fatalErr := output.FatalError(); fatalErr != nil {
				slogger.ErrorCtx(ctx, "failed to process message", "err", fatalErr)
				return fatalErr
			} else if nonFatalErr := output.NonFatalError(); nonFatalErr != nil {
				slogger.WarnCtx(ctx, "failed to process message", "err", nonFatalErr)
				return nonFatalErr
			} else {
				slogger.InfoCtx(ctx, "processed message", "output", output)
			}
			return nil
		},
	}
}

func provideConsumer(cfg consumer.Config, jobs job.Jobs) *consumer.Consumer {
	return either.Must(consumer.New(cfg, provideSQSClient(config.AWS()), func(jobType string) (sqsWorkerJob.Job, error) {
		return job.Get(jobType, jobs)
	}))
}

//lint:ignore U1000 used by wire
var consumerSet = wire.NewSet(
	provideConsumerConfig,
	provideConsumer,
)
