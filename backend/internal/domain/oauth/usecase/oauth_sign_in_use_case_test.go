package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/domain/oauth/model"
	"mickamy.com/sampay/internal/domain/oauth/usecase"
)

func TestCreateDirectUploadURL_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)

	// act
	sut := di.InitOAuthUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).OAuthSignIn
	got, err := sut.Do(ctx, usecase.OAuthSignInInput{
		Provider: model.OAuthProviderGoogle,
	})

	// assert
	require.NoError(t, err)
	assert.NotEmpty(t, got.AuthenticationURL)
}
