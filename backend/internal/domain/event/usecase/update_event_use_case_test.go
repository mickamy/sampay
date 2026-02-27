package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/mickamy/sampay/internal/domain/event/fixture"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/query"
	"github.com/mickamy/sampay/internal/domain/event/usecase"
	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func TestUpdateEvent_Do(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		ev := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		sut := usecase.NewUpdateEvent(infra)
		out, err := sut.Do(ctx, usecase.UpdateEventInput{
			ID:          ev.ID,
			Title:       "更新後タイトル",
			Description: "更新後説明",
			TotalAmount: 50000,
			TierCount:   5,
			HeldAt:      time.Now().Add(48 * time.Hour),
		})

		require.NoError(t, err)
		assert.Equal(t, "更新後タイトル", out.Event.Title)
		assert.Equal(t, 50000, out.Event.TotalAmount)
		assert.Equal(t, 5, out.Event.TierCount)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		sut := usecase.NewUpdateEvent(infra)
		_, err := sut.Do(ctx, usecase.UpdateEventInput{
			ID:          "nonexistent",
			Title:       "タイトル",
			TotalAmount: 10000,
			TierCount:   1,
			HeldAt:      time.Now(),
		})

		require.Error(t, err)
		require.ErrorIs(t, err, usecase.ErrUpdateEventNotFound)
	})

	t.Run("forbidden (other user's event)", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		owner := tseed.EndUser(t, infra.WriterDB)
		other := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), other.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		ev := fixture.Event(func(e *model.Event) { e.UserID = owner.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		sut := usecase.NewUpdateEvent(infra)
		_, err := sut.Do(ctx, usecase.UpdateEventInput{
			ID:          ev.ID,
			Title:       "タイトル",
			TotalAmount: 10000,
			TierCount:   1,
			HeldAt:      time.Now(),
		})

		require.Error(t, err)
		require.ErrorIs(t, err, usecase.ErrUpdateEventForbidden)
	})
}
