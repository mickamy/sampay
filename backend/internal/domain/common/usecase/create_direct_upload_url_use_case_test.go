package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/domain/common/fixture"
	"mickamy.com/sampay/internal/domain/common/usecase"
)

func TestCreateDirectUploadURL_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)

	// act
	sut := di.InitCommonUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).CreateDirectUploadURL
	got, err := sut.Do(ctx, usecase.CreateDirectUploadURLInput{
		S3Object: fixture.S3Object(nil),
	})

	// assert
	require.NoError(t, err)
	assert.NotEmpty(t, got.URL)
}
