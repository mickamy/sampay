package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/domain/registration/usecase"
)

func TestListUsageCategories_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)

	// act
	sut := di.InitRegistrationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).ListUsageCategories
	got, err := sut.Do(ctx, usecase.ListUsageCategoriesInput{})

	// assert
	require.NoError(t, err)
	assert.Len(t, got.Categories, 11)
}
