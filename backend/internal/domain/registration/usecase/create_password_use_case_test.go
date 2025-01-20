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
	"mickamy.com/sampay/internal/domain/registration/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/random"
)

func TestCreatePassword_Do(t *testing.T) {
	t.Parallel()

	token := either.Must(random.NewString(32))

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, db *database.Writer)
		assert  func(t *testing.T, got usecase.CreatePasswordOutput, err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer) {
				verification := authFixture.EmailVerificationVerified(func(m *authModel.EmailVerification) {
					m.IntentType = authModel.EmailVerificationIntentTypeSignUp
					m.Verified.Token = token
				})
				require.NoError(t, db.WithContext(ctx).Create(&verification).Error)
			},
			assert: func(t *testing.T, got usecase.CreatePasswordOutput, err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Session)
			},
		},
		{
			name: "invalid token",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer) {
				verification := authFixture.EmailVerificationVerified(func(m *authModel.EmailVerification) {
					m.IntentType = authModel.EmailVerificationIntentTypeSignUp
					m.Verified.Token = either.Must(random.NewString(32))
				})
				require.NoError(t, db.WithContext(ctx).Create(&verification).Error)
			},
			assert: func(t *testing.T, got usecase.CreatePasswordOutput, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "email verification not found")
				assert.Empty(t, got)
			},
		},
		{
			name: "already consumed",
			arrange: func(t *testing.T, ctx context.Context, db *database.Writer) {
				verification := authFixture.EmailVerificationConsumed(func(m *authModel.EmailVerification) {
					m.IntentType = authModel.EmailVerificationIntentTypeSignUp
					m.Verified.Token = token
				})
				require.NoError(t, db.WithContext(ctx).Create(&verification).Error)
			},
			assert: func(t *testing.T, got usecase.CreatePasswordOutput, err error) {
				require.Error(t, err)
				assert.ErrorContains(t, err, "email verification already consumed")
				assert.Empty(t, got)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := newReadWriter(t)
			ctx = contexts.SetAnonymousUserToken(ctx, token)
			tc.arrange(t, ctx, db.Writer())

			sut := di.InitRegistrationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).CreatePassword
			got, err := sut.Do(ctx, usecase.CreatePasswordInput{
				Password: gofakeit.Password(true, true, true, false, false, 12),
			})

			tc.assert(t, got, err)
		})
	}
}
