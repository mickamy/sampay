package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/either"
)

func TestAuthenticateUser_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	kvs := newKVS(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	session := either.Must(authModel.NewSession(user.ID))
	require.NoError(t, authRepository.NewSession(kvs).Create(ctx, session))

	// act
	sut := di.InitAuthUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), kvs).AuthenticateUser
	got, err := sut.Do(ctx, usecase.AuthenticateUserInput{
		Token: session.Tokens.Access.Value,
	})

	// assert
	require.NoError(t, err)
	assert.NotEmpty(t, got.User)
	assert.Equal(t, user.ID, got.User.ID)
	assert.Equal(t, user.Slug, got.User.Slug)
}
