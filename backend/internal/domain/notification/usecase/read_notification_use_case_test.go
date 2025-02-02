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
)

func TestReadNotification_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
	m := fixture.Notification(func(m *model.Notification) {
		m.UserID = user.ID
	})
	require.NoError(t, db.Writer().WithContext(ctx).Create(&m).Error)

	// act
	sut := di.InitNotificationUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).ReadNotification
	got, err := sut.Do(ctx, usecase.ReadNotificationInput{
		ID: m.ID,
	})

	// assert
	require.NoError(t, err)
	assert.Empty(t, got)
	var updated model.Notification
	require.NoError(t, db.Reader().WithContext(ctx).Joins("ReadStatus").First(&updated, "id = ?", m.ID).Error)
	assert.NotEmpty(t, updated.ReadStatus.ReadAt)
}
