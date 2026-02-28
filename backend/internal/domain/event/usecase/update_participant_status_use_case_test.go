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

func TestUpdateParticipantStatus_Do(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
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

		sut := usecase.NewUpdateParticipantStatus(infra)
		out, err := sut.Do(ctx, usecase.UpdateParticipantStatusInput{
			EventID:       ev.ID,
			ParticipantID: p.ID,
			Status:        model.ParticipantStatusConfirmed,
		})

		require.NoError(t, err)
		assert.Equal(t, model.ParticipantStatusConfirmed, out.Participant.Status)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		sut := usecase.NewUpdateParticipantStatus(infra)
		_, err := sut.Do(ctx, usecase.UpdateParticipantStatusInput{
			EventID:       "nonexistent",
			ParticipantID: "nonexistent",
			Status:        model.ParticipantStatusConfirmed,
		})

		require.ErrorIs(t, err, usecase.ErrUpdateParticipantStatusNotFound)
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

		p := fixture.EventParticipant(func(p *model.EventParticipant) {
			p.EventID = ev.ID
			p.Status = model.ParticipantStatusClaimed
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		sut := usecase.NewUpdateParticipantStatus(infra)
		_, err := sut.Do(ctx, usecase.UpdateParticipantStatusInput{
			EventID:       ev.ID,
			ParticipantID: p.ID,
			Status:        model.ParticipantStatusConfirmed,
		})

		require.ErrorIs(t, err, usecase.ErrUpdateParticipantStatusForbidden)
	})
}
