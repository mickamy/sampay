package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/ptr"
)

func TestUserProfile_Create(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	m := userFixture.UserProfile(func(m *model.UserProfile) {
		m.UserID = user.ID
	})

	// act
	sut := repository.NewUserProfile(db.WriterDB())
	err := sut.Create(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.UserProfile
	require.NoError(t, db.ReaderDB().WithContext(ctx).Where("user_id = ?", user.ID).First(&got).Error)
	assert.Equal(t, m.UserID, got.UserID)
	assert.Equal(t, m.Name, got.Name)
	assert.Equal(t, m.Bio, got.Bio)
	assert.Equal(t, m.ImageID, got.ImageID)
	assert.WithinDuration(t, m.CreatedAt, got.CreatedAt, time.Second)
	assert.WithinDuration(t, m.UpdatedAt, got.UpdatedAt, time.Second)
}

func TestUserProfile_Find(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(ctx context.Context, db *database.Writer, user *model.User)
		assert  func(t *testing.T, got *model.UserProfile, err error)
	}{
		{
			name: "found",
			arrange: func(ctx context.Context, db *database.Writer, user *model.User) {
				m := userFixture.UserProfile(func(m *model.UserProfile) {
					m.UserID = user.ID
				})
				require.NoError(t, db.WriterDB().WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got *model.UserProfile, err error) {
				require.NoError(t, err)
				require.NotNil(t, got)
			},
		},
		{
			name:    "not found",
			arrange: func(ctx context.Context, db *database.Writer, user *model.User) {},
			assert: func(t *testing.T, got *model.UserProfile, err error) {
				require.NoError(t, err)
				require.Nil(t, got)
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
			require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
			tc.arrange(ctx, db.Writer(), &user)

			// act
			sut := repository.NewUserProfile(db.ReaderDB())
			got, err := sut.FindBySlug(ctx, user.Slug)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestUserProfile_Update(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	m := userFixture.UserProfile(func(m *model.UserProfile) {
		m.UserID = user.ID
		m.Bio = nil
		m.SetImage(ptr.Of(commonFixture.S3Object(nil)))
	})
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&m).Error)
	m.Name = "new name"
	m.Bio = ptr.Of("new bio")
	m.SetImage(nil)

	// act
	sut := repository.NewUserProfile(db.WriterDB())
	err := sut.Update(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.UserProfile
	require.NoError(t, db.ReaderDB().WithContext(ctx).Where("user_id = ?", user.ID).First(&got).Error)
	assert.Equal(t, m.Name, got.Name)
	assert.Equal(t, *m.Bio, *got.Bio)
	assert.Nil(t, got.ImageID)
}

func TestUserProfile_Upsert(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(ctx context.Context, db *database.Writer, user *model.User) model.UserProfile
		assert  func(t *testing.T, got model.UserProfile, err error)
	}{
		{
			name: "not exists",
			arrange: func(ctx context.Context, db *database.Writer, user *model.User) model.UserProfile {
				return userFixture.UserProfile(func(m *model.UserProfile) {
					m.UserID = user.ID
					m.Name = "updated name"
					m.Bio = ptr.Of("updated bio")
				})
			},
			assert: func(t *testing.T, got model.UserProfile, err error) {
				require.NoError(t, err)
				assert.Equal(t, "updated name", got.Name)
				assert.Equal(t, "updated bio", *got.Bio)
			},
		},
		{
			name: "exists",
			arrange: func(ctx context.Context, db *database.Writer, user *model.User) model.UserProfile {
				m := userFixture.UserProfile(func(m *model.UserProfile) {
					m.UserID = user.ID
				})
				require.NoError(t, db.WriterDB().WithContext(ctx).Create(&m).Error)
				m.Name = "updated name"
				m.Bio = ptr.Of("updated bio")
				return m
			},
			assert: func(t *testing.T, got model.UserProfile, err error) {
				require.NoError(t, err)
				assert.Equal(t, "updated name", got.Name)
				assert.Equal(t, "updated bio", *got.Bio)
			},
		},
		{
			name: "exists with nil bio",
			arrange: func(ctx context.Context, db *database.Writer, user *model.User) model.UserProfile {
				m := userFixture.UserProfile(func(m *model.UserProfile) {
					m.UserID = user.ID
				})
				require.NoError(t, db.WriterDB().WithContext(ctx).Create(&m).Error)
				m.Name = "updated name"
				m.Bio = nil
				return m
			},
			assert: func(t *testing.T, got model.UserProfile, err error) {
				require.NoError(t, err)
				assert.Equal(t, "updated name", got.Name)
				assert.Nil(t, got.Bio)
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
			require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
			m := tc.arrange(ctx, db.Writer(), &user)

			// act
			sut := repository.NewUserProfile(db.WriterDB())
			err := sut.Upsert(ctx, &m)

			// assert
			require.NoError(t, err)
			var got model.UserProfile
			require.NoError(t, db.ReaderDB().WithContext(ctx).Where("user_id = ?", user.ID).First(&got).Error)
			tc.assert(t, got, err)
		})
	}
}
