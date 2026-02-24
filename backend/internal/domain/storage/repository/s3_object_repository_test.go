package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/storage/fixture"
	"github.com/mickamy/sampay/internal/domain/storage/query"
	"github.com/mickamy/sampay/internal/domain/storage/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

func TestS3Object_Create(t *testing.T) {
	t.Parallel()

	// arrange
	db := newReadWriter(t)
	m := fixture.S3Object(nil)

	// act
	sut := repository.NewS3Object(db.Writer.DB)
	err := sut.Create(t.Context(), &m)

	// assert
	require.NoError(t, err)
	got, err := query.S3Objects(db.Reader.DB).Where("id = ?", m.ID).First(t.Context())
	require.NoError(t, err)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.Bucket, got.Bucket)
	assert.Equal(t, m.Key, got.Key)
}

func TestS3Object_Get(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		// arrange
		db := newReadWriter(t)
		m := fixture.S3Object(nil)
		require.NoError(t, query.S3Objects(db.Writer.DB).Create(t.Context(), &m))

		// act
		sut := repository.NewS3Object(db.Reader.DB)
		got, err := sut.Get(t.Context(), m.ID)

		// assert
		require.NoError(t, err)
		assert.Equal(t, m.ID, got.ID)
		assert.Equal(t, m.Bucket, got.Bucket)
		assert.Equal(t, m.Key, got.Key)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		// arrange
		db := newReadWriter(t)

		// act
		sut := repository.NewS3Object(db.Reader.DB)
		_, err := sut.Get(t.Context(), "nonexistent")

		// assert
		assert.ErrorIs(t, err, database.ErrNotFound)
	})
}
