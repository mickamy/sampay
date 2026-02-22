package kvs

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/valkey-io/valkey-go"

	"github.com/mickamy/sampay/config"
)

var (
	kvsOnce    sync.Once
	kvsInst    *KVS
	kvsOpenErr error //nolint:errname // not a sentinel error; used for sync.Once result caching
)

func Open(cfg config.KVSConfig, opts ...Option) (*KVS, error) {
	kvsOnce.Do(func() {
		kvsInst, kvsOpenErr = New(cfg, opts...)
		if kvsOpenErr != nil {
			kvsOpenErr = fmt.Errorf("kvs: failed to open KVS: %w", kvsOpenErr)
		}
	})
	return kvsInst, kvsOpenErr
}

// KVS is a wrapper of valkey.Client
type KVS struct {
	client valkey.Client
}

type options struct {
	disableCache bool
}

type Option func(*options)

func WithDisableCache() Option {
	return func(o *options) {
		o.disableCache = true
	}
}

func New(cfg config.KVSConfig, opts ...Option) (*KVS, error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress:  []string{cfg.Address()},
		Username:     cfg.Username,
		Password:     cfg.Password,
		DisableCache: o.disableCache,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create valkey client: %w", err)
	}

	return &KVS{
		client: client,
	}, nil
}

type Storable interface {
	~string | ~[]byte | ~int64 | ~float64 | ~bool
}

func storableString[S Storable](v S) string {
	switch val := any(v).(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return fmt.Sprintf("%f", val)
	case bool:
		if val {
			return "1"
		}
		return "0"
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
}

func Get[T Storable](ctx context.Context, kvs *KVS, key string) (T, error) {
	res := kvs.client.Do(ctx, kvs.client.B().Get().Key(key).Build())

	var zero T

	if err := res.Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return zero, ErrKeyNotFound
		}
		return zero, fmt.Errorf("failed to get key: %w", err)
	}

	switch any(zero).(type) {
	case string:
		s, err := res.ToString()
		if err != nil {
			return zero, fmt.Errorf("failed to convert to string: %w", err)
		}
		v, ok := any(s).(T)
		if !ok {
			return zero, fmt.Errorf("unexpected type assertion failure: string to %T", zero)
		}
		return v, nil

	case []byte:
		b, err := res.AsBytes()
		if err != nil {
			return zero, fmt.Errorf("failed to convert to bytes: %w", err)
		}
		v, ok := any(b).(T)
		if !ok {
			return zero, fmt.Errorf("unexpected type assertion failure: bytes to %T", zero)
		}
		return v, nil

	case int64:
		i, err := res.ToInt64()
		if err != nil {
			return zero, fmt.Errorf("failed to convert to int64: %w", err)
		}
		v, ok := any(i).(T)
		if !ok {
			return zero, fmt.Errorf("unexpected type assertion failure: int64 to %T", zero)
		}
		return v, nil

	case float64:
		f, err := res.ToFloat64()
		if err != nil {
			return zero, fmt.Errorf("failed to convert to float64: %w", err)
		}
		v, ok := any(f).(T)
		if !ok {
			return zero, fmt.Errorf("unexpected type assertion failure: float64 to %T", zero)
		}
		return v, nil

	case bool:
		b, err := res.ToBool()
		if err != nil {
			return zero, fmt.Errorf("failed to convert to bool: %w", err)
		}
		v, ok := any(b).(T)
		if !ok {
			return zero, fmt.Errorf("unexpected type assertion failure: bool to %T", zero)
		}
		return v, nil
	}

	return zero, fmt.Errorf("unsupported type: %T", zero)
}

type Marshaler[T any, U Storable] interface {
	Marshal(v T) (U, error)
	Unmarshal(data U) (T, error)
}

func Memoize[T Storable, U any](
	ctx context.Context,
	kvs *KVS,
	key string,
	ttl time.Duration,
	exec func() (U, error),
	marshaler Marshaler[U, T],
) (U, error) {
	var zero U

	execAndStore := func() (U, error) {
		res, err := exec()
		if err != nil {
			return zero, err
		}

		storable, err := marshaler.Marshal(res)
		if err != nil {
			return res, errors.Join(ErrFailedToMarshalOnMemoize, fmt.Errorf("failed to marshal %T: %w", res, err))
		}

		cmd := kvs.client.B().Set().Key(key).Value(storableString(storable)).Px(ttl).Build()
		if err := kvs.client.Do(ctx, cmd).Error(); err != nil {
			return res, errors.Join(ErrFailedToSetOnMemoize, fmt.Errorf("failed to set key %s: %w", key, err))
		}

		return res, nil
	}

	cached, err := Get[T](ctx, kvs, key)
	if err != nil && !errors.Is(err, ErrKeyNotFound) {
		return zero, err
	}
	if err == nil {
		res, err := marshaler.Unmarshal(cached)
		if err != nil {
			return zero, errors.Join(
				ErrFailedToUnmarshalOnMemoize,
				fmt.Errorf("failed to unmarshal cached value for key %s: %w", key, err),
			)
		}
		return res, nil
	}

	return execAndStore()
}

func (c *KVS) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	if err := c.client.Do(ctx, c.client.B().Set().Key(key).Value(value).Px(exp).Build()).Error(); err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

func (c *KVS) Del(ctx context.Context, keys ...string) error {
	if err := c.client.Do(ctx, c.client.B().Del().Key(keys...).Build()).Error(); err != nil {
		return fmt.Errorf("failed to delete keys: %w", err)
	}
	return nil
}

func (c *KVS) Exists(ctx context.Context, keys ...string) (bool, error) {
	n, err := c.client.Do(ctx, c.client.B().Exists().Key(keys...).Build()).AsInt64()
	if err != nil {
		return false, fmt.Errorf("failed to check existence: %w", err)
	}
	return n > 0, nil
}

func (c *KVS) Ping(ctx context.Context) error {
	if err := c.client.Do(ctx, c.client.B().Ping().Build()).Error(); err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}
	return nil
}

func (c *KVS) Close() {
	c.client.Close()
}
