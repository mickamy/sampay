package di

import (
	"context"
	"fmt"

	"github.com/mickamy/go-sqs-worker/producer"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/infra/aws/s3"
	"github.com/mickamy/sampay/internal/infra/aws/sqs"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/infra/storage/kvs"
)

type Infra struct {
	_ context.Context  `inject:"param"` //nolint:containedctx // required by injector
	_ config.KVSConfig `inject:"provider:config.KVS"`
	// DB shares connection with WriterDB; do not close separately.
	DB       *database.DB       `inject:"provider:di.ProvideDB"`
	WriterDB *database.Writer   `inject:""`
	ReaderDB *database.Reader   `inject:""`
	KVS      *kvs.KVS           `inject:"provider:di.ProvideKVS"`
	S3       s3.Client          `inject:"provider:di.ProvideS3"`
	Producer *producer.Producer `inject:"provider:di.ProvideProducer"`
}

// Close releases all infrastructure resources.
// DB is not closed separately because it shares the underlying connection with WriterDB.
func (i *Infra) Close() error {
	i.KVS.Close()
	if err := i.WriterDB.Close(); err != nil {
		return fmt.Errorf("failed to close writer db: %w", err)
	}
	if err := i.ReaderDB.Close(); err != nil {
		return fmt.Errorf("failed to close reader db: %w", err)
	}

	return nil
}

func ProvideDB(
	ctx context.Context,
	commonCfg config.CommonConfig,
	databaseCfg config.DatabaseConfig,
) (*database.DB, error) {
	writer, err := ProvideWriterDB(ctx, commonCfg, databaseCfg)
	if err != nil {
		return nil, err
	}
	return writer.DB, nil
}

func ProvideWriterDB(
	ctx context.Context,
	commonCfg config.CommonConfig,
	databaseCfg config.DatabaseConfig,
) (*database.Writer, error) {
	db, err := database.Open(ctx, commonCfg, databaseCfg.WriterProvider(), database.RoleWriter)
	if err != nil {
		return nil, fmt.Errorf("failed to open writer db: %w", err)
	}
	return &database.Writer{DB: db}, nil
}

func ProvideReaderDB(
	ctx context.Context,
	commonCfg config.CommonConfig,
	databaseCfg config.DatabaseConfig,
) (*database.Reader, error) {
	db, err := database.Open(ctx, commonCfg, databaseCfg.ReaderProvider(), database.RoleReader)
	if err != nil {
		return nil, fmt.Errorf("failed to open reader db: %w", err)
	}
	return &database.Reader{DB: db}, nil
}

func ProvideS3(ctx context.Context) (s3.Client, error) {
	awsCfg := config.AWS()
	client, err := s3.New(ctx, awsCfg)
	if err != nil {
		return nil, fmt.Errorf("di: failed to initialize S3 client: %w", err)
	}
	return client, nil
}

func ProvideKVS(cfg config.KVSConfig) (*kvs.KVS, error) {
	kvStore, err := kvs.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("di: failed to initialize KVS: %w", err)
	}

	return kvStore, nil
}

func ProvideProducer(ctx context.Context, kvsCfg config.KVSConfig) (*producer.Producer, error) {
	awsCfg := config.AWS()
	sqsClient, err := sqs.New(ctx, awsCfg)
	if err != nil {
		return nil, fmt.Errorf("di: failed to initialize SQS client: %w", err)
	}

	redisURL := fmt.Sprintf("redis://:%s@%s", kvsCfg.Password, kvsCfg.Address())

	p, err := producer.New(producer.Config{
		WorkerQueueURL: awsCfg.SQSWorkerURL,
		RedisURL:       redisURL,
	}, sqsClient, nil)
	if err != nil {
		return nil, fmt.Errorf("di: failed to initialize producer: %w", err)
	}

	return p, nil
}
