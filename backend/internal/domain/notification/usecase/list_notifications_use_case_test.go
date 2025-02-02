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

func TestListNotification_Do(t *testing.T) {
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
	ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)

	// act
	sut := di.InitNotificationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).ListNotifications
	got, err := sut.Do(ctx, usecase.ListNotificationsInput{})

	// assert
	require.NoError(t, err)
	assert.Len(t, got.Notifications, 2)
}
