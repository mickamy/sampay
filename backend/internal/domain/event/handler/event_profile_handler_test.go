package handler_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/mickamy/contest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	eventv1 "github.com/mickamy/sampay/gen/event/v1"
	"github.com/mickamy/sampay/gen/event/v1/eventv1connect"
	"github.com/mickamy/sampay/internal/api/interceptor"
	"github.com/mickamy/sampay/internal/domain/event/fixture"
	"github.com/mickamy/sampay/internal/domain/event/handler"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/query"
	ufixture "github.com/mickamy/sampay/internal/domain/user/fixture"
	umodel "github.com/mickamy/sampay/internal/domain/user/model"
	uquery "github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/misc/i18n"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
	"github.com/mickamy/sampay/internal/test/ctest"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func TestEventProfile_GetEvent(t *testing.T) {
	t.Parallel()

	t.Run("returns event with user and payment methods without auth", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		pm := ufixture.UserPaymentMethod(func(m *umodel.UserPaymentMethod) { m.UserID = endUser.UserID })
		require.NoError(t, uquery.UserPaymentMethods(infra.WriterDB).Create(t.Context(), &pm))
		ev := fixture.Event(func(m *model.Event) { m.UserID = endUser.UserID; m.TierCount = 1 })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		tier := fixture.EventTier(func(m *model.EventTier) {
			m.EventID = ev.ID
			m.Tier = 1
			m.Count = 3
			m.Amount = ev.TotalAmount
		})
		require.NoError(t, query.EventTiers(infra.WriterDB).Create(t.Context(), &tier))
		p := fixture.EventParticipant(func(m *model.EventParticipant) {
			m.EventID = ev.ID
			m.Amount = 5000
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		// act
		var out eventv1.GetEventResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventProfileServiceHandler)(handler.NewEventProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventProfileServiceGetEventProcedure).
			In(&eventv1.GetEventRequest{Id: ev.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Equal(t, ev.ID, out.GetEvent().GetId())
		assert.Equal(t, endUser.UserID, out.GetUser().GetId())
		assert.Len(t, out.GetPaymentMethods(), 1)
		assert.Equal(t, pm.ID, out.GetPaymentMethods()[0].GetId())
		require.Len(t, out.GetParticipants(), 1)
		assert.Equal(t, p.ID, out.GetParticipants()[0].GetId())
		assert.Equal(t, int32(5000), out.GetParticipants()[0].GetAmount())
	})

	t.Run("returns not found for nonexistent event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventProfileServiceHandler)(handler.NewEventProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventProfileServiceGetEventProcedure).
			In(&eventv1.GetEventRequest{Id: "nonexistent"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusNotFound)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeNotFound, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorNotFound()), localized)
	})
}

func TestEventProfile_JoinEvent(t *testing.T) {
	t.Parallel()

	t.Run("joins event successfully", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		owner := tseed.EndUser(t, infra.WriterDB)
		ev := fixture.Event(func(m *model.Event) { m.UserID = owner.UserID; m.TierCount = 3; m.TotalAmount = 30000 })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		for i := 1; i <= 3; i++ {
			tier := fixture.EventTier(func(m *model.EventTier) { m.EventID = ev.ID; m.Tier = i; m.Count = 2 })
			require.NoError(t, query.EventTiers(infra.WriterDB).Create(t.Context(), &tier))
		}
		// recalc amounts
		ev.Tiers = []model.EventTier{
			{Tier: 1, Count: 2}, {Tier: 2, Count: 2}, {Tier: 3, Count: 2},
		}
		ev.CalcTierAmounts()
		for i := range ev.Tiers {
			_, err := infra.WriterDB.ExecContext(t.Context(),
				"UPDATE event_tiers SET amount = $1 WHERE event_id = $2 AND tier = $3",
				ev.Tiers[i].Amount, ev.ID, ev.Tiers[i].Tier,
			)
			require.NoError(t, err)
		}

		// act
		var out eventv1.JoinEventResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventProfileServiceHandler)(handler.NewEventProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventProfileServiceJoinEventProcedure).
			In(&eventv1.JoinEventRequest{
				EventId: ev.ID,
				Name:    "Taro Tanaka",
				Tier:    2,
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.NotEmpty(t, out.GetParticipant().GetId())
		assert.Equal(t, "Taro Tanaka", out.GetParticipant().GetName())
		assert.Equal(t, int32(2), out.GetParticipant().GetTier())
		assert.Equal(t, eventv1.ParticipantStatus_PARTICIPANT_STATUS_UNPAID, out.GetParticipant().GetStatus())
		assert.Positive(t, out.GetParticipant().GetAmount())
	})

	t.Run("returns error for empty name", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		owner := tseed.EndUser(t, infra.WriterDB)
		ev := fixture.Event(func(m *model.Event) { m.UserID = owner.UserID; m.TierCount = 1 })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		tier := fixture.EventTier(func(m *model.EventTier) {
			m.EventID = ev.ID
			m.Tier = 1
			m.Count = 3
			m.Amount = ev.TotalAmount
		})
		require.NoError(t, query.EventTiers(infra.WriterDB).Create(t.Context(), &tier))

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventProfileServiceHandler)(handler.NewEventProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventProfileServiceJoinEventProcedure).
			In(&eventv1.JoinEventRequest{
				EventId: ev.ID,
				Name:    "",
				Tier:    1,
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusBadRequest)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeInvalidArgument, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorNameRequired()), localized)
	})

	t.Run("returns error for invalid tier", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		owner := tseed.EndUser(t, infra.WriterDB)
		ev := fixture.Event(func(m *model.Event) { m.UserID = owner.UserID; m.TierCount = 3 })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		for i := 1; i <= 3; i++ {
			tier := fixture.EventTier(func(m *model.EventTier) { m.EventID = ev.ID; m.Tier = i; m.Count = 2; m.Amount = 5000 })
			require.NoError(t, query.EventTiers(infra.WriterDB).Create(t.Context(), &tier))
		}

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventProfileServiceHandler)(handler.NewEventProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventProfileServiceJoinEventProcedure).
			In(&eventv1.JoinEventRequest{
				EventId: ev.ID,
				Name:    "Taro Tanaka",
				Tier:    5, // invalid: tier_count is 3
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusBadRequest)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeInvalidArgument, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorInvalidTier()), localized)
	})

	t.Run("returns not found for nonexistent event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventProfileServiceHandler)(handler.NewEventProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventProfileServiceJoinEventProcedure).
			In(&eventv1.JoinEventRequest{
				EventId: "nonexistent",
				Name:    "Taro Tanaka",
				Tier:    1,
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusNotFound)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeNotFound, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorNotFound()), localized)
	})
}

func TestEventProfile_ClaimPayment(t *testing.T) {
	t.Parallel()

	t.Run("claims payment successfully", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		owner := tseed.EndUser(t, infra.WriterDB)
		ev := fixture.Event(func(m *model.Event) { m.UserID = owner.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		p := fixture.EventParticipant(func(m *model.EventParticipant) {
			m.EventID = ev.ID
			m.Status = model.ParticipantStatusUnpaid
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		// act
		var out eventv1.ClaimPaymentResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventProfileServiceHandler)(handler.NewEventProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventProfileServiceClaimPaymentProcedure).
			In(&eventv1.ClaimPaymentRequest{ParticipantId: p.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Equal(t, eventv1.ParticipantStatus_PARTICIPANT_STATUS_CLAIMED, out.GetParticipant().GetStatus())
	})

	t.Run("returns not found for nonexistent participant", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventProfileServiceHandler)(handler.NewEventProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventProfileServiceClaimPaymentProcedure).
			In(&eventv1.ClaimPaymentRequest{ParticipantId: "nonexistent"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusNotFound)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeNotFound, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorParticipantNotFound()), localized)
	})

	t.Run("returns error when already claimed", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		owner := tseed.EndUser(t, infra.WriterDB)
		ev := fixture.Event(func(m *model.Event) { m.UserID = owner.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		p := fixture.EventParticipant(func(m *model.EventParticipant) {
			m.EventID = ev.ID
			m.Status = model.ParticipantStatusClaimed
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventProfileServiceHandler)(handler.NewEventProfile(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventProfileServiceClaimPaymentProcedure).
			In(&eventv1.ClaimPaymentRequest{ParticipantId: p.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusBadRequest)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeFailedPrecondition, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorAlreadyClaimed()), localized)
	})
}
