package usecase_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
)

func TestRequestEmailVerification_Do(t *testing.T) {
	t.Parallel()

	email := gofakeit.GlobalFaker.Email()

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
				require.Empty(t, got)
			},
		},
		{
			name: "no verification exists",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
			},
			assert: func(t *testing.T, got usecase.RequestEmailVerificationOutput, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, got.Token)
				require.NotZero(t, got.ExpiresAt)
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
			sut := di.InitRegistrationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).RequestEmailVerification
			got, err := sut.Do(ctx, usecase.RequestEmailVerificationInput{
				Email: email,
			})

			// assert
			tc.assert(t, got, err)
		})
	}
}
