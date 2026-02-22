package kvs_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/infra/storage/kvs"
	"github.com/mickamy/sampay/internal/test/itest"
)

func TestGet_String(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	k := itest.NewKVS(t)

	err := k.Set(ctx, "key1", "hello", time.Minute)
	require.NoError(t, err)

	got, err := kvs.Get[string](ctx, k, "key1")
	require.NoError(t, err)
	assert.Equal(t, "hello", got)
}

func TestGet_Bytes(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	k := itest.NewKVS(t)

	err := k.Set(ctx, "key1", "bytes", time.Minute)
	require.NoError(t, err)

	got, err := kvs.Get[[]byte](ctx, k, "key1")
	require.NoError(t, err)
	assert.Equal(t, []byte("bytes"), got)
}

func TestGet_KeyNotFound(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	k := itest.NewKVS(t)

	_, err := kvs.Get[string](ctx, k, "nonexistent")
	assert.ErrorIs(t, err, kvs.ErrKeyNotFound)
}

type testData struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type jsonMarshaler struct{}

func (jsonMarshaler) Marshal(v testData) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonMarshaler) Unmarshal(b []byte) (testData, error) {
	var v testData
	err := json.Unmarshal(b, &v)
	return v, err
}

func TestMemoize_CacheMiss(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	k := itest.NewKVS(t)

	called := 0
	exec := func() (testData, error) {
		called++
		return testData{Name: "test", Value: 123}, nil
	}

	got, err := kvs.Memoize(ctx, k, "memo1", time.Minute, exec, jsonMarshaler{})
	require.NoError(t, err)
	assert.Equal(t, testData{Name: "test", Value: 123}, got)
	assert.Equal(t, 1, called)

	// verify value is stored in Valkey
	exists, err := k.Exists(ctx, "memo1")
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestMemoize_CacheHit(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	k := itest.NewKVS(t)

	called := 0
	exec := func() (testData, error) {
		called++
		return testData{Name: "test", Value: 123}, nil
	}

	// first call: cache miss, exec is called
	_, err := kvs.Memoize(ctx, k, "memo2", time.Minute, exec, jsonMarshaler{})
	require.NoError(t, err)
	assert.Equal(t, 1, called)

	// second call: cache hit, exec is not called
	got, err := kvs.Memoize(ctx, k, "memo2", time.Minute, exec, jsonMarshaler{})
	require.NoError(t, err)
	assert.Equal(t, testData{Name: "test", Value: 123}, got)
	assert.Equal(t, 1, called)
}

func TestMemoize_ExecError(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	k := itest.NewKVS(t)

	execErr := errors.New("exec failed")
	exec := func() (testData, error) {
		return testData{}, execErr
	}

	_, err := kvs.Memoize(ctx, k, "memo3", time.Minute, exec, jsonMarshaler{})
	assert.ErrorIs(t, err, execErr)

	// verify value is not stored in Valkey
	exists, err := k.Exists(ctx, "memo3")
	require.NoError(t, err)
	assert.False(t, exists)
}
