package usecase_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	registrationFixture "mickamy.com/sampay/internal/domain/registration/fixture"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
)

func TestCreatePassword_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
	verification := registrationFixture.EmailVerificationVerified(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&verification).Error)

	// act
	sut := di.InitRegistrationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).CreatePassword
	_, err := sut.Do(ctx, usecase.CreatePasswordInput{
		Email:    verification.Email,
		Password: gofakeit.Password(true, true, true, false, false, 12),
	})

	// assert
	require.NoError(t, err)
}
