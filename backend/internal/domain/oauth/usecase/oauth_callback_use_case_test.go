package usecase_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"mickamy.com/sampay/internal/di"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	oauthModel "mickamy.com/sampay/internal/domain/oauth/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	oauthLib "mickamy.com/sampay/internal/lib/oauth"
	"mickamy.com/sampay/internal/lib/oauth/mock_oauth"

	"mickamy.com/sampay/internal/domain/oauth/usecase"
)

func TestOauthCallback_Do(t *testing.T) {
	t.Parallel()

	t.Run("google", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctx := context.Background()
		db := newReadWriter(t)
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockGoogle := mock_oauth.NewMockGoogle(ctrl)
		code := "valid-code"
		mockGoogle.EXPECT().Validate(gomock.Eq(ctx), gomock.Eq(code)).Return(&oauthLib.Payload{
			Provider: oauthLib.ProviderGoogle,
			UID:      gofakeit.GlobalFaker.UUID(),
			Name:     gofakeit.GlobalFaker.Name(),
			Email:    gofakeit.GlobalFaker.Email(),
			Picture:  commonFixture.ImageURL(),
		}, nil)

		// act
		sut := usecase.NewOAuthCallback(
			mockGoogle,
			di.InitLibs().S3,
			db.Writer(),
			authRepository.NewAuthentication(db.WriterDB()),
			authRepository.NewEmailVerification(db.WriterDB()),
			authRepository.NewSession(newKVS(t)),
			userRepository.NewUser(db.WriterDB()),
			userRepository.NewUserProfile(db.WriterDB()),
		)
		got, err := sut.Do(ctx, usecase.OAuthCallbackInput{
			Provider: oauthModel.OAuthProviderGoogle,
			Code:     code,
		})

		// assert
		require.NoError(t, err)
		assert.NotEmpty(t, got.VerificationToken)
		assert.NotEmpty(t, got.Session.UserID)
		assert.NotEmpty(t, got.Session.Tokens.Access)
		assert.NotEmpty(t, got.Session.Tokens.Refresh)
	})
}
