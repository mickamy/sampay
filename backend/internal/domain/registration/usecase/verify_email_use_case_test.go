package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/di"
	registrationFixture "mickamy.com/sampay/internal/domain/registration/fixture"
	registrationModel "mickamy.com/sampay/internal/domain/registration/model"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/random"
)

func TestVerifyEmail_Do(t *testing.T) {
	t.Parallel()

	token := either.Must(random.NewString(32))

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, db *database.DB)
		assert  func(t *testing.T, got usecase.VerifyEmailOutput, err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := registrationFixture.EmailVerificationRequested(func(m *registrationModel.EmailVerification) {
					m.Requested.Token = token
				})
				assert.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got usecase.VerifyEmailOutput, err error) {
				require.NoError(t, err)
				assert.Empty(t, got)
			},
		},
		{
			name: "token expired",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := registrationFixture.EmailVerificationRequested(func(m *registrationModel.EmailVerification) {
					m.Requested.Token = token
					m.Requested.ExpiresAt = time.Now().Add(-time.Second)
				})
				assert.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got usecase.VerifyEmailOutput, err error) {
				assert.ErrorIs(t, err, registrationModel.ErrEmailVerificationTokenExpired)
				assert.ErrorContains(t, err, "failed to verify email verification")
				assert.Empty(t, got)
			},
		},
		{
			name: "token not expired",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := registrationFixture.EmailVerificationRequested(func(m *registrationModel.EmailVerification) {
					m.Requested.Token = token
					m.Requested.ExpiresAt = time.Now().Add(5 * time.Second)
				})
				assert.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got usecase.VerifyEmailOutput, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "token verified",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := registrationFixture.EmailVerificationVerified(func(m *registrationModel.EmailVerification) {
					m.Requested.Token = token
				})
				assert.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got usecase.VerifyEmailOutput, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "token consumed",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := registrationFixture.EmailVerificationConsumed(func(m *registrationModel.EmailVerification) {
					m.Requested.Token = token
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
			sut := di.InitRegistrationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).VerifyEmail
			got, err := sut.Do(ctx, usecase.VerifyEmailInput{Token: token})

			// assert
			tc.assert(t, got, err)
		})
	}
}
