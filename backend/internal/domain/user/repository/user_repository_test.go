package repository_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/repository"
)

func TestUser_Create(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	user := userFixture.User(nil)
	db := newReadWriter(t)

	// act
	sut := repository.NewUser(db.WriterDB())
	err := sut.Create(ctx, &user)

	// assert
	require.NoError(t, err)
	var got model.User
	err = db.ReaderDB().WithContext(ctx).First(&got, "id = ?", user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Slug, got.Slug)
}

func TestUser_Find(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	user := userFixture.User(nil)
	db := newReadWriter(t)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)

	// act
	sut := repository.NewUser(db.WriterDB())
	got, err := sut.Find(ctx, user.ID)

	// assert
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Slug, got.Slug)
}

func TestUser_FindBySlug(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	user := userFixture.User(nil)
	db := newReadWriter(t)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)

	// act
	sut := repository.NewUser(db.WriterDB())
	got, err := sut.FindBySlug(ctx, user.Slug)

	// assert
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Slug, got.Slug)
}

func TestUser_FindByEmail(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	user := userFixture.User(nil)
	db := newReadWriter(t)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)

	// act
	sut := repository.NewUser(db.WriterDB())
	got, err := sut.FindByEmail(ctx, user.Email)

	// assert
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Slug, got.Slug)
}

func TestUser_FindByEmailOrSlug(t *testing.T) {
	t.Parallel()

	email := gofakeit.GlobalFaker.Email()
	slug := gofakeit.GlobalFaker.Username()

	tcs := []struct {
		name        string
		emailOrSlug string
	}{
		{
			name:        "email",
			emailOrSlug: email,
		},
		{
			name:        "slug",
			emailOrSlug: slug,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctx := context.Background()
			user := userFixture.User(func(m *model.User) {
				m.Email = email
				m.Slug = slug
			})
			db := newReadWriter(t)
			require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)

			// act
			sut := repository.NewUser(db.WriterDB())
			got, err := sut.FindByEmailOrSlug(ctx, tc.emailOrSlug)

			// assert
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, user.ID, got.ID)
			assert.Equal(t, user.Slug, got.Slug)
		})
	}
}

func TestUser_Update(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	user := userFixture.User(nil)
	db := newReadWriter(t)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	user.Slug = gofakeit.GlobalFaker.Username()

	// act
	sut := repository.NewUser(db.WriterDB())
	err := sut.Update(ctx, &user)

	// assert
	require.NoError(t, err)
	var got model.User
	err = db.ReaderDB().WithContext(ctx).First(&got, "id = ?", user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, user.ID, got.ID)
	assert.Equal(t, user.Slug, got.Slug)
}
