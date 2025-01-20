package usecase_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/random"
)

func TestResetPassword_Do(t *testing.T) {
	t.Parallel()

	token := either.Must(random.NewString(32))

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, db *database.Writer, userID string)
		assert  func(t *testing.T, got usecase.ResetPasswordOutput, err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer, userID string) {
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&auth).Error)
				verification := authFixture.EmailVerificationVerified(func(m *authModel.EmailVerification) {
					m.IntentType = authModel.EmailVerificationIntentTypeSignUp
					m.Email = auth.Identifier
					m.Verified.Token = token
				})
				require.NoError(t, db.WithContext(ctx).Create(&verification).Error)
			},
			assert: func(t *testing.T, got usecase.ResetPasswordOutput, err error) {
				require.NoError(t, err)
				assert.Empty(t, got)
			},
		},
		{
			name: "invalid token",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer, userID string) {
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&auth).Error)
				verification := authFixture.EmailVerificationVerified(func(m *authModel.EmailVerification) {
					m.IntentType = authModel.EmailVerificationIntentTypeSignUp
					m.Email = auth.Identifier
					m.Verified.Token = either.Must(random.NewString(32))
				})
				require.NoError(t, db.WithContext(ctx).Create(&verification).Error)
			},
			assert: func(t *testing.T, got usecase.ResetPasswordOutput, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "email verification not found")
			},
		},
		{
			name: "already consumed",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer, userID string) {
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
				})
				require.NoError(t, db.WithContext(ctx).Create(&auth).Error)
				verification := authFixture.EmailVerificationConsumed(func(m *authModel.EmailVerification) {
					m.IntentType = authModel.EmailVerificationIntentTypeSignUp
					m.Email = auth.Identifier
					m.Verified.Token = token
				})
				require.NoError(t, db.WithContext(ctx).Create(&verification).Error)
			},
			assert: func(t *testing.T, got usecase.ResetPasswordOutput, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "email verification already consumed")
			},
		},
		{
			name: "authentication not found",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer, userID string) {
				verification := authFixture.EmailVerificationVerified(func(m *authModel.EmailVerification) {
					m.IntentType = authModel.EmailVerificationIntentTypeSignUp
					m.Verified.Token = token
				})
				require.NoError(t, db.WithContext(ctx).Create(&verification).Error)
			},
			assert: func(t *testing.T, got usecase.ResetPasswordOutput, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "authentication not found")
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctx := context.Background()
			db := newReadWriter(t)
			user := userFixture.User(nil)
			require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
			ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
			tc.arrange(t, ctx, db.Writer(), user.ID)

			// act
			sut := di.InitAuthUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).ResetPassword
			got, err := sut.Do(ctx, usecase.ResetPasswordInput{
				Token:    token,
				Password: gofakeit.Password(true, true, true, false, false, 12),
			})

			// assert
			tc.assert(t, got, err)
		})
	}
}
