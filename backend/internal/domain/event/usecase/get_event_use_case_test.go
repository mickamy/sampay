package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/event/fixture"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/query"
	"github.com/mickamy/sampay/internal/domain/event/usecase"
	ufixture "github.com/mickamy/sampay/internal/domain/user/fixture"
	umodel "github.com/mickamy/sampay/internal/domain/user/model"
	uquery "github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func TestGetEvent_Do(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)

		pm := ufixture.UserPaymentMethod(func(m *umodel.UserPaymentMethod) {
			m.UserID = endUser.UserID
		})
		require.NoError(t, uquery.UserPaymentMethods(infra.WriterDB).Create(t.Context(), &pm))

		ev := fixture.Event(func(e *model.Event) {
			e.UserID = endUser.UserID
			e.TotalAmount = 10000
			e.TierCount = 3
		})
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		p1 := fixture.EventParticipant(func(p *model.EventParticipant) {
			p.EventID = ev.ID
			p.Tier = 3
		})
		p2 := fixture.EventParticipant(func(p *model.EventParticipant) {
			p.EventID = ev.ID
			p.Tier = 1
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p1))
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p2))

		sut := usecase.NewGetEvent(infra)
		out, err := sut.Do(t.Context(), usecase.GetEventInput{ID: ev.ID})

		require.NoError(t, err)
		assert.Equal(t, ev.ID, out.Event.ID)
		assert.Equal(t, endUser.UserID, out.User.UserID)
		assert.Len(t, out.PaymentMethods, 1)
		require.Len(t, out.Event.Participants, 2)
		// totalWeight = 3 + 1 = 4
		amountByID := make(map[string]int, len(out.Event.Participants))
		for _, p := range out.Event.Participants {
			amountByID[p.ID] = p.Amount
		}
		assert.Equal(t, 10000*3/4, amountByID[p1.ID])
		assert.Equal(t, 10000*1/4, amountByID[p2.ID])
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)

		sut := usecase.NewGetEvent(infra)
		_, err := sut.Do(t.Context(), usecase.GetEventInput{ID: "nonexistent"})

		require.Error(t, err)
		require.ErrorIs(t, err, usecase.ErrGetEventNotFound)
	})
}
