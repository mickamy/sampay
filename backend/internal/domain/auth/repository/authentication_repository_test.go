package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/auth/repository"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/slices"
)

func TestAuthentication_Create(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
		m.UserID = user.ID
	})

	// act
	sut := repository.NewAuthentication(db.WriterDB())
	err := sut.Create(ctx, &auth)

	// assert
	require.NoError(t, err)
	var got authModel.Authentication
	err = db.ReaderDB().WithContext(ctx).First(&got, "id = ?", auth.ID).Error
	require.NoError(t, err)
	assert.Equal(t, auth.ID, got.ID)
	assert.Equal(t, auth.UserID, got.UserID)
	assert.Equal(t, auth.Type, got.Type)
	assert.Equal(t, auth.Identifier, got.Identifier)
	assert.Equal(t, auth.Secret, got.Secret)
	assert.Equal(t, auth.MFAEnabled, got.MFAEnabled)
	assert.WithinDuration(t, auth.CreatedAt, got.CreatedAt, time.Second)
	assert.WithinDuration(t, auth.UpdatedAt, got.UpdatedAt, time.Second)
}

func TestAuthentication_FindByKey(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
		m.UserID = user.ID
	})
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&auth).Error)

	// act
	sut := repository.NewAuthentication(db.WriterDB())
	got, err := sut.FindByKey(ctx, repository.AuthenticationKey{
		UserID:     auth.UserID,
		Type:       auth.Type,
		Identifier: auth.Identifier,
	})

	// assert
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, auth.ID, got.ID)
	assert.Equal(t, auth.UserID, got.UserID)
	assert.Equal(t, auth.Type, got.Type)
	assert.Equal(t, auth.Identifier, got.Identifier)
	assert.Equal(t, auth.Secret, got.Secret)
	assert.Equal(t, auth.MFAEnabled, got.MFAEnabled)
}

func TestAuthentication_ListByUserID(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	var auths []authModel.Authentication
	for i := 0; i < 3; i++ {
		auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
			m.UserID = user.ID
		})
		require.NoError(t, db.WriterDB().WithContext(ctx).Create(&auth).Error)
		auths = append(auths, auth)
	}

	// act
	sut := repository.NewAuthentication(db.WriterDB())
	gots, err := sut.ListByUserID(ctx, user.ID)

	// assert
	require.NoError(t, err)
	require.Len(t, gots, len(auths))
	for i := range gots {
		got, found := slices.Find(auths, func(authentication authModel.Authentication) bool {
			return authentication.ID == gots[i].ID
		})
		require.True(t, found)
		assert.Equal(t, got.ID, gots[i].ID)
		assert.Equal(t, got.UserID, gots[i].UserID)
		assert.Equal(t, got.Type, gots[i].Type)
		assert.Equal(t, got.Identifier, gots[i].Identifier)
		assert.Equal(t, got.Secret, gots[i].Secret)
		assert.Equal(t, got.MFAEnabled, gots[i].MFAEnabled)
		assert.WithinDuration(t, got.CreatedAt, gots[i].CreatedAt, time.Second)
		assert.WithinDuration(t, got.UpdatedAt, gots[i].UpdatedAt, time.Second)
	}
}

func TestAuthentication_Update(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
		m.UserID = user.ID
	})
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&auth).Error)
	auth.Secret = "new-secret"

	// act
	sut := repository.NewAuthentication(db.WriterDB())
	err := sut.Update(ctx, &auth)

	// assert
	require.NoError(t, err)
	var got authModel.Authentication
	err = db.ReaderDB().WithContext(ctx).First(&got, "id = ?", auth.ID).Error
	require.NoError(t, err)
	assert.Equal(t, auth.ID, got.ID)
	assert.Equal(t, auth.UserID, got.UserID)
	assert.Equal(t, auth.Type, got.Type)
	assert.Equal(t, auth.Identifier, got.Identifier)
	assert.Equal(t, auth.Secret, got.Secret)
	assert.Equal(t, auth.MFAEnabled, got.MFAEnabled)
}
