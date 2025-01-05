package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/test/infra"
)

func TestCreateSessionUseCase_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := NewReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
		m.UserID = user.ID
	})
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&auth).Error)

	// act
	sut := di.InitAuthUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), infra.NewKVS(t)).CreateSession
	got, err := sut.Do(ctx, usecase.CreateSessionInput{
		Email:    auth.Identifier,
		Password: "P@ssw0rd",
	})
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.Session.UserID)
	assert.NotEmpty(t, got.Session.Tokens.Access.Value)
	assert.NotEmpty(t, got.Session.Tokens.Refresh.Value)
	assert.NotEmpty(t, got.Session.Tokens.Access.ExpiresAt)
	assert.NotEmpty(t, got.Session.Tokens.Refresh.ExpiresAt)
}
