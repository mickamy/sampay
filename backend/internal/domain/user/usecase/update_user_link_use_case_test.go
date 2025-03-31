package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	"mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/domain/user/usecase"
	"mickamy.com/sampay/internal/infra/storage/database"
	"mickamy.com/sampay/internal/lib/ptr"
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

	// act
	sut := di.InitUserUseCase(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).UpdateUserLink
	got, err := sut.Do(ctx, usecase.UpdateUserLinkInput{
		ID:           m.ID,
		ProviderType: ptr.Of(model.UserLinkProviderTypeOther),
		URI:          ptr.Of("https://example.com"),
		Name:         ptr.Of("updated"),
		QRImage:      ptr.Of(commonFixture.S3Object(nil)),
	})

	// assert
	require.NoError(t, err)
	assert.Empty(t, got)
	var updated model.UserLink
	require.NoError(t, db.ReaderDB().WithContext(ctx).Scopes(database.Scope(repository.UserLinkJoinDisplayAttribute).Gorm()).First(&updated, "id = ?", m.ID).Error)
	assert.Equal(t, m.ID, updated.ID)
	assert.Equal(t, model.UserLinkProviderTypeOther, updated.ProviderType)
	assert.Equal(t, "https://example.com", updated.URI)
	assert.Equal(t, "updated", updated.DisplayAttribute.Name)
	assert.NotEmpty(t, updated.QRCodeID)
}
