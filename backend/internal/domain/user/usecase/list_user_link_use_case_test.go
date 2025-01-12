package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/usecase"
)

func TestListUserLink_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := fixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	m := fixture.UserLink(func(m *model.UserLink) {
		m.UserID = user.ID
	})
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&m).Error)

	// act
	sut := di.InitUserUseCase(db.ReaderDB(), db, db.Writer(), db.Reader(), newKVS(t)).ListUserLink
	got, err := sut.Do(ctx, usecase.ListUserLinkInput{UserID: user.ID})

	// assert
	require.NoError(t, err)
	assert.Len(t, got.Links, 1)
	assert.NotEmpty(t, got.Links[0])
}
