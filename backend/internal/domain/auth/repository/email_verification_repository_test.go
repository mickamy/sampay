package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/auth/fixture"
	"mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/auth/repository"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/random"
)

func TestEmailVerification_Create(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	m := fixture.EmailVerification(nil)

	// act
	sut := repository.NewEmailVerification(db.WriterDB())
	err := sut.Create(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.EmailVerification
	require.NoError(t, db.WriterDB().WithContext(ctx).First(&got, "id = ?", m.ID).Error)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.Email, got.Email)
	assert.WithinDuration(t, m.CreatedAt, got.CreatedAt, time.Second)
}

func TestEmailVerification_FindByEmail(t *testing.T) {
	t.Parallel()

	email := gofakeit.GlobalFaker.Email()

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, db *database.DB)
		assert  func(t *testing.T, got *model.EmailVerification, err error)
	}{
		{
			name: "found",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := fixture.EmailVerification(func(m *model.EmailVerification) {
					m.Email = email
				})
				require.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got *model.EmailVerification, err error) {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.NotEmpty(t, got.ID)
				assert.NotEmpty(t, got.Email)
				assert.NotEmpty(t, got.CreatedAt)
			},
		}, {
			name:    "not found",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {},
			assert: func(t *testing.T, got *model.EmailVerification, err error) {
				require.NoError(t, err)
				assert.Nil(t, got)
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
			sut := repository.NewEmailVerification(db.WriterDB())
			got, err := sut.FindByEmail(ctx, email)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestEmailVerification_FindByRequestedTokenAndPinCode(t *testing.T) {
	t.Parallel()

	token := either.Must(random.NewString(32))
	pin := either.Must(random.NewPinCode(6))

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, db *database.DB)
		assert  func(t *testing.T, got *model.EmailVerification, err error)
	}{
		{
			name: "found",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := fixture.EmailVerificationRequested(func(m *model.EmailVerification) {
					m.Requested.Token = token
					m.Requested.PINCode = pin
				})
				require.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got *model.EmailVerification, err error) {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.NotEmpty(t, got.ID)
				assert.NotEmpty(t, got.Email)
				assert.NotEmpty(t, got.CreatedAt)
			},
		}, {
			name: "not found (pin code is different)",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := fixture.EmailVerification(func(m *model.EmailVerification) {
					m.Requested.Token = token
				})
				require.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got *model.EmailVerification, err error) {
				require.NoError(t, err)
				assert.Nil(t, got)
			},
		}, {
			name: "not found (token is different)",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := fixture.EmailVerification(func(m *model.EmailVerification) {
					m.Requested.PINCode = pin
				})
				require.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got *model.EmailVerification, err error) {
				require.NoError(t, err)
				assert.Nil(t, got)
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
			sut := repository.NewEmailVerification(db.WriterDB())
			got, err := sut.FindByRequestedTokenAndPinCode(ctx, token, pin)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestEmailVerification_FindByVerifiedToken(t *testing.T) {
	t.Parallel()

	token := either.Must(random.NewString(32))

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, db *database.DB)
		assert  func(t *testing.T, got *model.EmailVerification, err error)
	}{
		{
			name: "found",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := fixture.EmailVerificationVerified(func(m *model.EmailVerification) {
					m.Verified.Token = token
				})
				require.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got *model.EmailVerification, err error) {
				require.NoError(t, err)
				require.NotNil(t, got)
				assert.NotEmpty(t, got.ID)
				assert.NotEmpty(t, got.Email)
				assert.NotEmpty(t, got.CreatedAt)
			},
		}, {
			name: "not found",
			arrange: func(t *testing.T, ctx context.Context, db *database.DB) {
				m := fixture.EmailVerification(nil)
				require.NoError(t, db.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got *model.EmailVerification, err error) {
				require.NoError(t, err)
				assert.Nil(t, got)
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
			sut := repository.NewEmailVerification(db.WriterDB())
			got, err := sut.FindByVerifiedToken(ctx, token)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestEmailVerification_Update(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	m := fixture.EmailVerification(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&m).Error)
	m.Email = gofakeit.GlobalFaker.Email()

	// act
	sut := repository.NewEmailVerification(db.WriterDB())
	err := sut.Update(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.EmailVerification
	require.NoError(t, db.WriterDB().WithContext(ctx).First(&got, "id = ?", m.ID).Error)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.Email, got.Email)
}
