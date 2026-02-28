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
			Title:       "updated title",
			Description: "updated description",
			TotalAmount: 50000,
			TierCount:   5,
			HeldAt:      time.Now().Add(48 * time.Hour),
			Tiers: []usecase.TierConfig{
				{Tier: 1, Count: 2},
				{Tier: 2, Count: 2},
				{Tier: 3, Count: 1},
				{Tier: 4, Count: 1},
				{Tier: 5, Count: 1},
			},
		})

		require.NoError(t, err)
		assert.Equal(t, "updated title", out.Event.Title)
		assert.Equal(t, 50000, out.Event.TotalAmount)
		assert.Equal(t, 5, out.Event.TierCount)
		require.Len(t, out.Event.Tiers, 5)
	})

	t.Run("locked (participant claimed)", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		ev := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		p := fixture.EventParticipant(func(p *model.EventParticipant) {
			p.EventID = ev.ID
			p.Status = model.ParticipantStatusClaimed
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		sut := usecase.NewUpdateEvent(infra)
		_, err := sut.Do(ctx, usecase.UpdateEventInput{
			ID:          ev.ID,
			Title:       "title",
			TotalAmount: 10000,
			TierCount:   1,
			Tiers:       []usecase.TierConfig{{Tier: 1, Count: 5}},
			HeldAt:      time.Now(),
		})

		require.Error(t, err)
		require.ErrorIs(t, err, usecase.ErrUpdateEventLocked)
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
			Title:       "title",
			TotalAmount: 10000,
			TierCount:   1,
			Tiers:       []usecase.TierConfig{{Tier: 1, Count: 5}},
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
			Title:       "title",
			TotalAmount: 10000,
			TierCount:   1,
			Tiers:       []usecase.TierConfig{{Tier: 1, Count: 5}},
			HeldAt:      time.Now(),
		})

		require.Error(t, err)
		require.ErrorIs(t, err, usecase.ErrUpdateEventForbidden)
	})
}
