package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
)

func TestGetUser_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(func(m *model.User) {
		m.Profile = userFixture.UserProfile(nil)
		m.Links = []model.UserLink{
			userFixture.UserLink(func(m *model.UserLink) {
				m.DisplayAttribute = userFixture.UserLinkDisplayAttribute(nil)
			}),
		}
	})
	require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
	ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)

	// act
	sut := di.InitUserUseCase(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).GetUser
	got, err := sut.Do(ctx, usecase.GetUserInput{Slug: user.Slug})

	// assert
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.NotEmpty(t, got.Profile)
	require.NotEmpty(t, got.Links)
	assert.NotEmpty(t, got.Links[0].DisplayAttribute)
}
