package usecase_test

import (
	"testing"

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

func TestListEventParticipants_Do(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		ev := fixture.Event(func(e *model.Event) {
			e.UserID = endUser.UserID
			e.TotalAmount = 12000
			e.TierCount = 3
		})
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		p1 := fixture.EventParticipant(func(p *model.EventParticipant) {
			p.EventID = ev.ID
			p.Tier = 3
			p.Amount = 9000
		})
		p2 := fixture.EventParticipant(func(p *model.EventParticipant) {
			p.EventID = ev.ID
			p.Tier = 1
			p.Amount = 3000
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p1))
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p2))

		sut := usecase.NewListEventParticipants(infra)
		out, err := sut.Do(ctx, usecase.ListEventParticipantsInput{EventID: ev.ID})

		require.NoError(t, err)
		require.Len(t, out.Participants, 2)
		amountByID := make(map[string]int, len(out.Participants))
		for _, p := range out.Participants {
			amountByID[p.ID] = p.Amount
		}
		assert.Equal(t, 9000, amountByID[p1.ID])
		assert.Equal(t, 3000, amountByID[p2.ID])
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		sut := usecase.NewListEventParticipants(infra)
		_, err := sut.Do(ctx, usecase.ListEventParticipantsInput{EventID: "nonexistent"})

		require.ErrorIs(t, err, usecase.ErrListEventParticipantsNotFound)
	})

	t.Run("forbidden", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		owner := tseed.EndUser(t, infra.WriterDB)
		other := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), other.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		ev := fixture.Event(func(e *model.Event) { e.UserID = owner.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		sut := usecase.NewListEventParticipants(infra)
		_, err := sut.Do(ctx, usecase.ListEventParticipantsInput{EventID: ev.ID})

		require.ErrorIs(t, err, usecase.ErrListEventParticipantsForbidden)
	})
}
