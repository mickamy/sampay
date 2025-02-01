package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/domain/message/fixture"
	"mickamy.com/sampay/internal/domain/message/model"
	"mickamy.com/sampay/internal/domain/message/repository"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
)

func TestMessage_Create(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
	m := fixture.Message(func(m *model.Message) {
		m.ReceiverID = user.ID
	})

	// act
	sut := repository.NewMessage(db.WriterDB())
	err := sut.Create(ctx, &m)

	// assert
	require.NoError(t, err)
	var got model.Message
	require.NoError(t, db.Reader().WithContext(ctx).First(&got, "id = ?", m.ID).Error)
	assert.Equal(t, m.ID, got.ID)
	assert.Equal(t, m.SenderName, got.SenderName)
	assert.Equal(t, m.ReceiverID, got.ReceiverID)
	assert.Equal(t, m.Content, got.Content)
	assert.WithinDuration(t, m.CreatedAt, got.CreatedAt, time.Second)
}

func TestMessage_ListByReceiverID(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	db := newReadWriter(t)
	user := userFixture.User(nil)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&user).Error)
	messages := []model.Message{
		fixture.Message(func(m *model.Message) {
			m.ReceiverID = user.ID
		}),
		fixture.Message(func(m *model.Message) {
			m.ReceiverID = user.ID
		}),
	}
	for i := range messages {
		require.NoError(t, db.Writer().WithContext(ctx).Create(&messages[i]).Error)
	}

	// act
	sut := repository.NewMessage(db.ReaderDB())
	got, err := sut.ListByReceiverID(ctx, user.ID)

	// assert
	require.NoError(t, err)
	assert.Len(t, got, 2)
	for i := range got {
		assert.Equal(t, messages[i].ID, got[i].ID)
		assert.Equal(t, messages[i].SenderName, got[i].SenderName)
		assert.Equal(t, messages[i].ReceiverID, got[i].ReceiverID)
		assert.Equal(t, messages[i].Content, got[i].Content)
		assert.WithinDuration(t, messages[i].CreatedAt, got[i].CreatedAt, time.Second)
	}
}
