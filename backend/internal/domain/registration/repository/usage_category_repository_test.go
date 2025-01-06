package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/domain/registration/model"
	"mickamy.com/sampay/internal/domain/registration/repository"
)

func TestUsageCategory_List(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)

	// act
	sut := repository.NewUsageCategory(db.WriterDB())
	got, err := sut.List(ctx)

	// assert
	require.NoError(t, err)
	assert.Len(t, got, 11)
}

func TestUsageCategory_Upsert(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	m := &model.UsageCategory{CategoryType: "other", DisplayOrder: 100}

	// act
	sut := repository.NewUsageCategory(db.WriterDB())
	err := sut.Upsert(ctx, m)

	// assert
	require.NoError(t, err)
	var got model.UsageCategory
	require.NoError(t, db.WriterDB().WithContext(ctx).Where("category_type = ?", m.CategoryType).First(&got).Error)
	assert.Equal(t, m, &got)
}
