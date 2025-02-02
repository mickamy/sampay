package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/notification/fixture"
	"mickamy.com/sampay/internal/domain/notification/model"
	"mickamy.com/sampay/internal/domain/notification/repository"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/ulid"
)

func TestNotification_Create(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
	m := fixture.Notification(func(m *model.Notification) {
		m.UserID = user.ID
	})

	// act
	sut := repository.NewNotification(db.WriterDB())
	err := sut.Create(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.Notification
	require.NoError(t, db.Reader().WithContext(ctx).First(&got, "id = ?", m.ID).Error)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.UserID, got.UserID)
	assert.Equal(t, m.Subject, got.Subject)
	assert.Equal(t, m.Body, got.Body)
}

func TestNotification_Find(t *testing.T) {
	t.Parallel()

	id := ulid.New()

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, writer *database.Writer, userID string)
		assert  func(t *testing.T, got *model.Notification, err error)
	}{
		{
			name: "found",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer, userID string) {
				m := fixture.Notification(func(m *model.Notification) {
					m.ID = id
					m.UserID = userID
				})
				require.NoError(t, writer.WithContext(ctx).Create(&m).Error)
			},
			assert: func(t *testing.T, got *model.Notification, err error) {
				require.NoError(t, err)
				assert.NotNil(t, got)
			},
		},
		{
			name:    "not found",
			arrange: func(t *testing.T, ctx context.Context, writer *database.Writer, userID string) {},
			assert: func(t *testing.T, got *model.Notification, err error) {
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
			user := userFixture.User(nil)
			require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
			tc.arrange(t, ctx, db.Writer(), user.ID)

			// act
			sut := repository.NewNotification(db.ReaderDB())
			got, err := sut.Find(ctx, id)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestNotification_ListByUserID(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
	m1 := fixture.Notification(func(m *model.Notification) {
		m.UserID = user.ID
	})
	m2 := fixture.Notification(func(m *model.Notification) {
		m.UserID = user.ID
	})
	require.NoError(t, db.Writer().WithContext(ctx).Create(&m1).Error)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&m2).Error)

	// act
	sut := repository.NewNotification(db.ReaderDB())
	got, err := sut.ListByUserID(ctx, user.ID)

	// assert
	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, m1.ID, got[0].ID)
	assert.Equal(t, m1.UserID, got[0].UserID)
	assert.Equal(t, m1.Subject, got[0].Subject)
	assert.Equal(t, m1.Body, got[0].Body)
	assert.Equal(t, m2.ID, got[1].ID)
	assert.Equal(t, m2.UserID, got[1].UserID)
	assert.Equal(t, m2.Subject, got[1].Subject)
	assert.Equal(t, m2.Body, got[1].Body)
}

func TestNotification_Update(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
	m := fixture.Notification(func(m *model.Notification) {
		m.UserID = user.ID
		require.NoError(t, db.Writer().WithContext(ctx).Create(&m).Error)
		m.Subject = "updated"
		m.Body = "updated"
	})

	// act
	sut := repository.NewNotification(db.WriterDB())
	err := sut.Update(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.Notification
	require.NoError(t, db.Reader().WithContext(ctx).First(&got, "id = ?", m.ID).Error)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.UserID, got.UserID)
	assert.Equal(t, m.Subject, got.Subject)
	assert.Equal(t, m.Body, got.Body)
}
