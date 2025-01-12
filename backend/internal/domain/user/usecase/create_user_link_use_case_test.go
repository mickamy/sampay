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
	"mickamy.com/sampay/internal/lib/contexts"
)

func TestCreateUserLink_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := fixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
	m := fixture.UserLink(func(m *model.UserLink) {
		m.UserID = user.ID
		m.DisplayAttribute = fixture.UserLinkDisplayAttribute(nil)
	})

	// act
	sut := di.InitUserUseCase(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).CreateUserLink
	output, err := sut.Do(ctx, usecase.CreateUserLinkInput{
		ProviderType: m.ProviderType,
		URI:          m.URI,
		Name:         m.DisplayAttribute.Name,
	})

	// assert
	require.NoError(t, err)
	assert.Empty(t, output)
}
