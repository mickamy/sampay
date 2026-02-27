package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/event/fixture"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/query"
	"github.com/mickamy/sampay/internal/domain/event/usecase"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func TestClaimPayment_Do(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)

		ev := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		p := fixture.EventParticipant(func(p *model.EventParticipant) {
			p.EventID = ev.ID
			p.Status = model.ParticipantStatusUnpaid
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		sut := usecase.NewClaimPayment(infra)
		out, err := sut.Do(t.Context(), usecase.ClaimPaymentInput{ParticipantID: p.ID})

		require.NoError(t, err)
		assert.Equal(t, model.ParticipantStatusClaimed, out.Participant.Status)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)

		sut := usecase.NewClaimPayment(infra)
		_, err := sut.Do(t.Context(), usecase.ClaimPaymentInput{ParticipantID: "nonexistent"})

		require.ErrorIs(t, err, usecase.ErrClaimPaymentNotFound)
	})

	t.Run("already claimed", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)

		ev := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		p := fixture.EventParticipant(func(p *model.EventParticipant) {
			p.EventID = ev.ID
			p.Status = model.ParticipantStatusClaimed
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		sut := usecase.NewClaimPayment(infra)
		_, err := sut.Do(t.Context(), usecase.ClaimPaymentInput{ParticipantID: p.ID})

		require.ErrorIs(t, err, usecase.ErrClaimPaymentAlreadyClaimed)
	})
}
