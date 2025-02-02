package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/domain/notification/fixture"
	"mickamy.com/sampay/internal/domain/notification/model"
	"mickamy.com/sampay/internal/domain/notification/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
)

func TestCountUnreadNotifications_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
	ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
	unread := fixture.Notification(func(m *model.Notification) {
		m.UserID = user.ID
	})
	require.NoError(t, db.Writer().WithContext(ctx).Create(&unread).Error)
	read := fixture.Notification(func(m *model.Notification) {
		m.UserID = user.ID
		m.ReadStatus = fixture.NotificationReadStatusRead(func(m *model.NotificationReadStatus) {
			m.UserID = user.ID
		})
	})
	require.NoError(t, db.Writer().WithContext(ctx).Create(&read).Error)

	// act
	sut := di.InitNotificationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).CountUnreadNotifications
	got, err := sut.Do(ctx, usecase.CountUnreadNotificationsInput{})

	// assert
	require.NoError(t, err)
	assert.Equal(t, 1, got.Count)
}
