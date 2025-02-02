package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	"mickamy.com/sampay/internal/domain/message/model"
	"mickamy.com/sampay/internal/domain/message/usecase"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/language"
)

func TestSendMessage_Do(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	ctx = contexts.SetLanguage(ctx, language.Japanese)
	db := newReadWriter(t)
	receiver := userFixture.User(nil)
	require.NoError(t, db.Writer().WithContext(ctx).Create(&receiver).Error)

	// act
	sut := di.InitMessageUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).SendMessage
	got, err := sut.Do(ctx, usecase.SendMessageInput{
		SenderName:   "sender",
		ReceiverSlug: receiver.Slug,
		Content:      "content",
	})

	// assert
	require.NoError(t, err)
	assert.Empty(t, got)
	var created model.Message
	require.NoError(t, db.Reader().WithContext(ctx).Last(&created).Error)
	assert.Equal(t, "sender", created.SenderName)
	assert.Equal(t, receiver.ID, created.ReceiverID)
	assert.Equal(t, "content", created.Content)
}
