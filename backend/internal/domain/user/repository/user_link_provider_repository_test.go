package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/repository"
)

func TestUserLinkProvider_Upsert(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	m := fixture.UserLinkProvider(nil)

	// act
	sut := repository.NewUserLinkProvider(db.WriterDB())
	err := sut.Upsert(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.UserLinkProvider
	require.NoError(t, db.ReaderDB().WithContext(ctx).First(&got, m.Type).Error)
	require.Equal(t, m.DisplayOrder, got.DisplayOrder)
}
