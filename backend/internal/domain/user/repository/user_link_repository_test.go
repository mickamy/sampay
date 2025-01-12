package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/repository"
)

func TestUserLink_Create(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := fixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	m := fixture.UserLink(func(m *model.UserLink) {
		m.UserID = user.ID
	})

	// act
	sut := repository.NewUserLink(db.WriterDB())
	err := sut.Create(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.UserLink
	require.NoError(t, db.ReaderDB().WithContext(ctx).First(&got, "id = ?", m.ID).Error)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.UserID, got.UserID)
	assert.Equal(t, m.ProviderType, got.ProviderType)
	assert.Equal(t, m.URI, got.URI)
}

func TestUserLink_ListByUserID(t *testing.T) {
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
	sut := repository.NewUserLink(db.ReaderDB())
	got, err := sut.ListByUserID(ctx, user.ID)

	// assert
	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, m.ID, got[0].ID)
	assert.Equal(t, m.UserID, got[0].UserID)
	assert.Equal(t, m.ProviderType, got[0].ProviderType)
	assert.Equal(t, m.URI, got[0].URI)
}

func TestUserLink_Find(t *testing.T) {
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
	sut := repository.NewUserLink(db.ReaderDB())
	got, err := sut.Find(ctx, m.ID)

	// assert
	require.NoError(t, err)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.UserID, got.UserID)
	assert.Equal(t, m.ProviderType, got.ProviderType)
	assert.Equal(t, m.URI, got.URI)
}

func TestUserLink_Update(t *testing.T) {
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
	sut := repository.NewUserLink(db.WriterDB())
	m.URI = "new-uri"
	err := sut.Update(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.UserLink
	require.NoError(t, db.ReaderDB().WithContext(ctx).First(&got, "id = ?", m.ID).Error)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.UserID, got.UserID)
	assert.Equal(t, m.ProviderType, got.ProviderType)
	assert.Equal(t, m.URI, got.URI)
}

func TestUserLink_Delete(t *testing.T) {
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
	sut := repository.NewUserLink(db.WriterDB())
	err := sut.Delete(ctx, m.ID)

	// assert
	require.NoError(t, err)
	var got model.UserLink
	require.Error(t, db.ReaderDB().WithContext(ctx).First(&got, "id = ?", m.ID).Error)
}
