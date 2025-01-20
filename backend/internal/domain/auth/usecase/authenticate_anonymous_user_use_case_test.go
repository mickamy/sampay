package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	"mickamy.com/sampay/internal/domain/auth/usecase"
)

func TestAuthenticateAnonymousUser_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	m := authFixture.EmailVerificationVerified(nil)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&m).Error)

	// act
	sut := di.InitAuthUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).AuthenticateAnonymousUser
	got, err := sut.Do(ctx, usecase.AuthenticateAnonymousUserInput{
		Token: m.Verified.Token,
	})

	// assert
	require.NoError(t, err)
	assert.Empty(t, got)
}
