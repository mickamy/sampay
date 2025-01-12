package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/domain/user/usecase"
)

func TestUpdateUserLink_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := fixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	m := fixture.UserLink(func(m *model.UserLink) {
		m.UserID = user.ID
		m.DisplayAttribute = fixture.UserLinkDisplayAttribute(nil)
	})
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&m).Error)

	m.ProviderType = model.UserLinkProviderTypeOther
	m.URI = "https://example.com"
	m.DisplayAttribute.Name = "example"
	m.DisplayAttribute.DisplayOrder = m.DisplayAttribute.DisplayOrder + 1

	// act
	sut := di.InitUserUseCase(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).UpdateUserLink
	got, err := sut.Do(ctx, usecase.UpdateUserLinkInput{
		ID:           m.ID,
		ProviderType: &m.ProviderType,
		URI:          &m.URI,
		Name:         &m.DisplayAttribute.Name,
		DisplayOrder: &m.DisplayAttribute.DisplayOrder,
	})

	// assert
	require.NoError(t, err)
	assert.Empty(t, got)
	var updated model.UserLink
	require.NoError(t, db.ReaderDB().WithContext(ctx).Scopes(database.Scope(repository.UserLinkJoinDisplayAttribute).Gorm()).First(&updated, "id = ?", m.ID).Error)
	assert.Equal(t, m.ID, updated.ID)
	assert.Equal(t, m.ProviderType, updated.ProviderType)
	assert.Equal(t, m.URI, updated.URI)
	assert.Equal(t, m.DisplayAttribute.Name, updated.DisplayAttribute.Name)
	assert.Equal(t, m.DisplayAttribute.DisplayOrder, updated.DisplayAttribute.DisplayOrder)
}
