package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/random"
)

func TestVerifyEmail_Do(t *testing.T) {
	t.Parallel()

	email := gofakeit.GlobalFaker.Email()
	pin := either.Must(random.NewPinCode(6))

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, db *database.DB)
		assert  func(t *testing.T, got usecase.VerifyEmailOutput, err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := authFixture.EmailVerificationRequested(func(m *authModel.EmailVerification) {
					m.Email = email
					m.Requested.PINCode = pin
				})
				assert.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got usecase.VerifyEmailOutput, err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Session)
			},
		},
		{
			name: "different email",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := authFixture.EmailVerificationRequested(func(m *authModel.EmailVerification) {
					m.Email = email + "different"
					m.Requested.PINCode = pin
				})
				assert.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got usecase.VerifyEmailOutput, err error) {
				assert.ErrorIs(t, err, usecase.ErrVerifyEmailInvalidToken)
			},
		},
		{
			name: "pin code expired",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := authFixture.EmailVerificationRequested(func(m *authModel.EmailVerification) {
					m.Email = email
					m.Requested.PINCode = pin
					m.Requested.ExpiresAt = time.Now().Add(-time.Second)
				})
				assert.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got usecase.VerifyEmailOutput, err error) {
				assert.ErrorIs(t, err, authModel.ErrEmailVerificationTokenExpired)
				assert.ErrorContains(t, err, "failed to verify email verification")
				assert.Empty(t, got)
			},
		},
		{
			name: "pin code not expired",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := authFixture.EmailVerificationRequested(func(m *authModel.EmailVerification) {
					m.Email = email
					m.Requested.PINCode = pin
					m.Requested.ExpiresAt = time.Now().Add(5 * time.Second)
				})
				assert.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got usecase.VerifyEmailOutput, err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Session)
			},
		},
		{
			name: "pin code verified",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := authFixture.EmailVerificationVerified(func(m *authModel.EmailVerification) {
					m.Email = email
					m.Requested.PINCode = pin
				})
				assert.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got usecase.VerifyEmailOutput, err error) {
				require.ErrorIs(t, err, usecase.ErrVerifyEmailInvalidToken)
				assert.Empty(t, got)
			},
		},
		{
			name: "pin code consumed",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := authFixture.EmailVerificationConsumed(func(m *authModel.EmailVerification) {
					m.Email = email
					m.Requested.PINCode = pin
				})
				assert.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got usecase.VerifyEmailOutput, err error) {
				assert.ErrorIs(t, err, usecase.ErrVerifyEmailInvalidToken)
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
			tc.arrange(t, ctx, db.WriterDB())

			// act
			sut := di.InitAuthUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).VerifyEmail
			got, err := sut.Do(ctx, usecase.VerifyEmailInput{Email: email, PINCode: pin})

			// assert
			tc.assert(t, got, err)
		})
	}
}
