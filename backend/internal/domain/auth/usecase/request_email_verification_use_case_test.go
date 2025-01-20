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
)

func TestRequestEmailVerification_Do(t *testing.T) {
	t.Parallel()

	email := gofakeit.GlobalFaker.Email()

	t.Run("sign_up", func(t *testing.T) {
		t.Parallel()

		tcs := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, db *database.DB)
			assert  func(t *testing.T, got usecase.RequestEmailVerificationOutput, err error)
		}{
			{
				name: "email already exists",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
					user := userFixture.User(nil)
					require.NoError(t, db.WithContext(ctx).Create(&user).Error)
					auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
						m.Identifier = email
						m.UserID = user.ID
					})
					require.NoError(t, db.WithContext(ctx).Create(&auth).Error)
				},
				assert: func(t *testing.T, got usecase.RequestEmailVerificationOutput, err error) {
					require.ErrorIs(t, err, usecase.ErrRequestEmailVerificationEmailAlreadyExists)
					assert.Empty(t, got)
				},
			},
			{
				name: "no verification exists",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				},
				assert: func(t *testing.T, got usecase.RequestEmailVerificationOutput, err error) {
					require.NoError(t, err)
					assert.NotEmpty(t, got.Token)
				},
			},
		}

		for _, tc := range tcs {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				// arrange
				ctx := context.Background()
				ctx = contexts.SetLanguage(ctx, "ja")
				db := newReadWriter(t)
				tc.arrange(t, ctx, db.WriterDB())

				// act
				sut := di.InitAuthUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).RequestEmailVerification
				got, err := sut.Do(ctx, usecase.RequestEmailVerificationInput{
					IntentType: authModel.EmailVerificationIntentTypeSignUp,
					Email:      email,
				})

				// assert
				tc.assert(t, got, err)
			})
		}
	})

	t.Run("reset_password", func(t *testing.T) {
		tcs := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, db *database.DB)
			assert  func(t *testing.T, got usecase.RequestEmailVerificationOutput, err error)
		}{
			{
				name: "email already exists",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
					user := userFixture.User(nil)
					require.NoError(t, db.WithContext(ctx).Create(&user).Error)
					auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
						m.Identifier = email
						m.UserID = user.ID
					})
					require.NoError(t, db.WithContext(ctx).Create(&auth).Error)
				},
				assert: func(t *testing.T, got usecase.RequestEmailVerificationOutput, err error) {
					require.NoError(t, err)
					assert.NotEmpty(t, got.Token)
				},
			},
			{
				name: "no authentication exists",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
					user := userFixture.User(nil)
					require.NoError(t, db.WithContext(ctx).Create(&user).Error)
				},
				assert: func(t *testing.T, got usecase.RequestEmailVerificationOutput, err error) {
					require.Error(t, err)
					require.ErrorContains(t, err, "authentication not found")
					assert.Empty(t, got.Token)
				},
			},
			{
				name: "no verification exists",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
					user := userFixture.User(nil)
					require.NoError(t, db.WithContext(ctx).Create(&user).Error)
					auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
						m.Identifier = email
						m.UserID = user.ID
					})
					require.NoError(t, db.WithContext(ctx).Create(&auth).Error)
				},
				assert: func(t *testing.T, got usecase.RequestEmailVerificationOutput, err error) {
					require.NoError(t, err)
					assert.NotEmpty(t, got.Token)
				},
			},
		}

		for _, tc := range tcs {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				// arrange
				ctx := context.Background()
				ctx = contexts.SetLanguage(ctx, "ja")
				db := newReadWriter(t)
				tc.arrange(t, ctx, db.WriterDB())

				// act
				sut := di.InitAuthUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).RequestEmailVerification
				got, err := sut.Do(ctx, usecase.RequestEmailVerificationInput{
					IntentType: authModel.EmailVerificationIntentTypeResetPassword,
					Email:      email,
				})

				// assert
				tc.assert(t, got, err)
			})
		}
	})
}
