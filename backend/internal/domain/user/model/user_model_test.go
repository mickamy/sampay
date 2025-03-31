package model_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/infra/storage/database"
)

func TestUser_BeforeCreate(t *testing.T) {
	t.Parallel()

	t.Run("slug", func(t *testing.T) {
		t.Parallel()

		slug := gofakeit.GlobalFaker.Username()

		tcs := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, writer *database.Writer)
			assert  func(t *testing.T, err error)
		}{
			{
				name: "taken",
				arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) {
					m := fixture.User(func(m *model.User) {
						m.Slug = slug
					})
					require.NoError(t, writer.WithContext(ctx).Create(&m).Error)
				},
				assert: func(t *testing.T, err error) {
					require.Error(t, err)
					assert.ErrorIs(t, err, model.ErrUserSlugAlreadyTaken)
				},
			},
			{
				name: "not taken",
				arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) {
				},
				assert: func(t *testing.T, err error) {
					require.NoError(t, err)
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

				tc.arrange(t, ctx, db.Writer())

				// act
				m := fixture.User(func(m *model.User) {
					m.Slug = slug
				})
				err := db.Writer().WithContext(ctx).Create(&m).Error

				// assert
				tc.assert(t, err)
			})
		}
	})

	t.Run("email", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.GlobalFaker.Email()

		tcs := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, writer *database.Writer)
			assert  func(t *testing.T, err error)
		}{
			{
				name: "taken",
				arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) {
					m := fixture.User(func(m *model.User) {
						m.Email = email
					})
					require.NoError(t, writer.WithContext(ctx).Create(&m).Error)
				},
				assert: func(t *testing.T, err error) {
					require.Error(t, err)
					assert.ErrorIs(t, err, model.ErrUserEmailAlreadyTaken)
				},
			},
			{
				name: "not taken",
				arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) {
				},
				assert: func(t *testing.T, err error) {
					require.NoError(t, err)
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

				tc.arrange(t, ctx, db.Writer())

				// act
				m := fixture.User(func(m *model.User) {
					m.Email = email
				})
				err := db.Writer().WithContext(ctx).Create(&m).Error

				// assert
				tc.assert(t, err)
			})
		}
	})
}

func TestUser_BeforeUpdate(t *testing.T) {
	t.Parallel()

	t.Run("slug", func(t *testing.T) {
		t.Parallel()

		slug := gofakeit.GlobalFaker.Username()

		tcs := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, writer *database.Writer)
			assert  func(t *testing.T, err error)
		}{
			{
				name: "taken",
				arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) {
					m := fixture.User(func(m *model.User) {
						m.Slug = slug
					})
					require.NoError(t, writer.WithContext(ctx).Create(&m).Error)
				},
				assert: func(t *testing.T, err error) {
					require.Error(t, err)
					assert.ErrorIs(t, err, model.ErrUserSlugAlreadyTaken)
				},
			},
			{
				name: "not taken",
				arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) {
				},
				assert: func(t *testing.T, err error) {
					require.NoError(t, err)
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

				tc.arrange(t, ctx, db.Writer())

				// act
				m := fixture.User(nil)
				require.NoError(t, db.Writer().WithContext(ctx).Create(&m).Error)
				m.Slug = slug
				err := db.Writer().WithContext(ctx).Updates(&m).Error

				// assert
				tc.assert(t, err)
			})
		}
	})

	t.Run("email", func(t *testing.T) {
		t.Parallel()

		email := gofakeit.GlobalFaker.Email()

		tcs := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, writer *database.Writer)
			assert  func(t *testing.T, err error)
		}{
			{
				name: "taken",
				arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) {
					m := fixture.User(func(m *model.User) {
						m.Email = email
					})
					require.NoError(t, writer.WithContext(ctx).Create(&m).Error)
				},
				assert: func(t *testing.T, err error) {
					require.Error(t, err)
					assert.ErrorIs(t, err, model.ErrUserEmailAlreadyTaken)
				},
			},
			{
				name: "not taken",
				arrange: func(t *testing.T, ctx context.Context, writer *database.Writer) {
				},
				assert: func(t *testing.T, err error) {
					require.NoError(t, err)
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

				tc.arrange(t, ctx, db.Writer())

				// act
				m := fixture.User(nil)
				require.NoError(t, db.Writer().WithContext(ctx).Create(&m).Error)
				m.Email = email
				err := db.Writer().WithContext(ctx).Updates(&m).Error

				// assert
				tc.assert(t, err)
			})
		}
	})
}
