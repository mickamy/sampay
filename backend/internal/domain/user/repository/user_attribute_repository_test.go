package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/repository"
)

func TestUserAttribute_Create(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.WriterDB().WithContext(ctx).Create(&user).Error)
	m := userFixture.UserAttribute(func(m *model.UserAttribute) {
		m.UserID = user.ID
	})

	// act
	sut := repository.NewUserAttribute(db.WriterDB())
	err := sut.Create(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.UserAttribute
	require.NoError(t, db.ReaderDB().WithContext(ctx).Where("user_id = ?", user.ID).First(&got).Error)
	assert.Equal(t, m.UserID, got.UserID)
	assert.Equal(t, m.UsageCategoryType, got.UsageCategoryType)
	assert.WithinDuration(t, m.CreatedAt, got.CreatedAt, time.Second)
	assert.WithinDuration(t, m.UpdatedAt, got.UpdatedAt, time.Second)
}

func TestUserAttribute_Find(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		name    string
		arrange func(ctx context.Context, db *database.Writer, user *model.User)
		assert  func(t *testing.T, got *model.UserAttribute, err error)
	}{
		{
			name: "found",
			arrange: func(ctx context.Context, db *database.Writer, user *model.User) {
				m := userFixture.UserAttribute(func(m *model.UserAttribute) {
					m.UserID = user.ID
				})
				require.NoError(t, db.WriterDB().WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got *model.UserAttribute, err error) {
				require.NoError(t, err)
				require.NotNil(t, got)
			},
		},
		{
			name:    "not found",
			arrange: func(ctx context.Context, db *database.Writer, user *model.User) {},
			assert: func(t *testing.T, got *model.UserAttribute, err error) {
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
			sut := repository.NewUserAttribute(db.ReaderDB())
			got, err := sut.Find(ctx, user.ID)

			// assert
			tc.assert(t, got, err)
		})
	}
}
