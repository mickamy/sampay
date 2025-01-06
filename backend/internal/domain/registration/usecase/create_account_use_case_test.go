package usecase_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/domain/registration/usecase"
)

func TestCreateAccountUseCase_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)

	// act
	sut := di.InitRegistrationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).CreateAccount
	got, err := sut.Do(ctx, usecase.CreateAccountInput{
		Email:    gofakeit.GlobalFaker.Email(),
		Password: gofakeit.GlobalFaker.Password(true, true, true, true, false, 12),
	})

	// assert
	require.NoError(t, err)
	assert.NotEmpty(t, got.Session.UserID)
	assert.NotEmpty(t, got.Session.Tokens.Access.Value)
	assert.NotEmpty(t, got.Session.Tokens.Refresh.Value)
	assert.NotEmpty(t, got.Session.Tokens.Access.ExpiresAt)
	assert.NotEmpty(t, got.Session.Tokens.Refresh.ExpiresAt)
}
