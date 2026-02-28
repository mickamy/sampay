package handler_test

import (
	"net/http"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/mickamy/contest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	eventv1 "github.com/mickamy/sampay/gen/event/v1"
	"github.com/mickamy/sampay/gen/event/v1/eventv1connect"
	"github.com/mickamy/sampay/internal/api/interceptor"
	"github.com/mickamy/sampay/internal/domain/event/fixture"
	"github.com/mickamy/sampay/internal/domain/event/handler"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/query"
	"github.com/mickamy/sampay/internal/misc/i18n"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
	"github.com/mickamy/sampay/internal/test/ctest"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func TestEventService_ListMyEvents(t *testing.T) {
	t.Parallel()

	t.Run("returns events for authenticated user", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		ev := fixture.Event(func(m *model.Event) { m.UserID = userID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		// act
		var out eventv1.ListMyEventsResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceListMyEventsProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.ListMyEventsRequest{}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		require.Len(t, out.GetEvents(), 1)
		assert.Equal(t, ev.ID, out.GetEvents()[0].GetId())
	})

	t.Run("returns unauthorized without auth header", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceListMyEventsProcedure).
			In(&eventv1.ListMyEventsRequest{}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusUnauthorized)
	})
}

func TestEventService_CreateEvent(t *testing.T) {
	t.Parallel()

	t.Run("creates event successfully", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)
		heldAt := time.Now().Add(24 * time.Hour)

		// act
		var out eventv1.CreateEventResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceCreateEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.CreateEventRequest{
				Input: &eventv1.EventInput{
					Title:       "bounenkai",
					Description: "2025 bounenkai",
					TotalAmount: 30000,
					TierCount:   1,
					HeldAt:      timestamppb.New(heldAt),
					Tiers:       []*eventv1.TierConfig{{Tier: 1, Count: 5}},
				},
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.NotEmpty(t, out.GetEvent().GetId())
		assert.Equal(t, "bounenkai", out.GetEvent().GetTitle())
		assert.Equal(t, int32(30000), out.GetEvent().GetTotalAmount())
		assert.Equal(t, int32(1), out.GetEvent().GetTierCount())
	})

	t.Run("returns error for empty title", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceCreateEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.CreateEventRequest{
				Input: &eventv1.EventInput{
					Title:       "",
					TotalAmount: 30000,
					TierCount:   1,
					HeldAt:      timestamppb.New(time.Now().Add(24 * time.Hour)),
					Tiers:       []*eventv1.TierConfig{{Tier: 1, Count: 5}},
				},
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusBadRequest)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeInvalidArgument, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorTitleRequired()), localized)
	})

	t.Run("returns error for invalid tier count", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceCreateEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.CreateEventRequest{
				Input: &eventv1.EventInput{
					Title:       "bounenkai",
					TotalAmount: 30000,
					TierCount:   2,
					HeldAt:      timestamppb.New(time.Now().Add(24 * time.Hour)),
					Tiers:       []*eventv1.TierConfig{{Tier: 1, Count: 3}, {Tier: 2, Count: 2}},
				},
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusBadRequest)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeInvalidArgument, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorTierCountInvalid()), localized)
	})
}

func TestEventService_UpdateEvent(t *testing.T) {
	t.Parallel()

	t.Run("updates event successfully", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		ev := fixture.Event(func(m *model.Event) { m.UserID = userID; m.TierCount = 1 })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		tier := fixture.EventTier(func(m *model.EventTier) {
			m.EventID = ev.ID
			m.Tier = 1
			m.Count = 3
			m.Amount = ev.TotalAmount
		})
		require.NoError(t, query.EventTiers(infra.WriterDB).Create(t.Context(), &tier))

		// act
		var out eventv1.UpdateEventResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUpdateEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UpdateEventRequest{
				Id: ev.ID,
				Input: &eventv1.EventInput{
					Title:       "shinnenkai",
					Description: "updated",
					TotalAmount: 50000,
					TierCount:   1,
					HeldAt:      timestamppb.New(ev.HeldAt),
					Tiers:       []*eventv1.TierConfig{{Tier: 1, Count: 5}},
				},
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Equal(t, "shinnenkai", out.GetEvent().GetTitle())
		assert.Equal(t, int32(50000), out.GetEvent().GetTotalAmount())
	})

	t.Run("returns not found for nonexistent event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUpdateEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UpdateEventRequest{
				Id: "nonexistent",
				Input: &eventv1.EventInput{
					Title:       "shinnenkai",
					TotalAmount: 50000,
					TierCount:   1,
					HeldAt:      timestamppb.New(time.Now().Add(24 * time.Hour)),
					Tiers:       []*eventv1.TierConfig{{Tier: 1, Count: 5}},
				},
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusNotFound)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeNotFound, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorNotFound()), localized)
	})

	t.Run("returns forbidden for other user's event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)
		otherUser := tseed.EndUser(t, infra.WriterDB)
		ev := fixture.Event(func(m *model.Event) { m.UserID = otherUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		tier := fixture.EventTier(func(m *model.EventTier) {
			m.EventID = ev.ID
			m.Tier = 1
			m.Count = 1
			m.Amount = ev.TotalAmount
		})
		require.NoError(t, query.EventTiers(infra.WriterDB).Create(t.Context(), &tier))

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUpdateEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UpdateEventRequest{
				Id: ev.ID,
				Input: &eventv1.EventInput{
					Title:       "shinnenkai",
					TotalAmount: 50000,
					TierCount:   1,
					HeldAt:      timestamppb.New(time.Now().Add(24 * time.Hour)),
					Tiers:       []*eventv1.TierConfig{{Tier: 1, Count: 5}},
				},
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusForbidden)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodePermissionDenied, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorForbidden()), localized)
	})

	t.Run("returns locked when participant has claimed", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		ev := fixture.Event(func(m *model.Event) { m.UserID = userID; m.TierCount = 1 })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		tier := fixture.EventTier(func(m *model.EventTier) {
			m.EventID = ev.ID
			m.Tier = 1
			m.Count = 1
			m.Amount = ev.TotalAmount
		})
		require.NoError(t, query.EventTiers(infra.WriterDB).Create(t.Context(), &tier))
		p := fixture.EventParticipant(func(m *model.EventParticipant) {
			m.EventID = ev.ID
			m.Status = model.ParticipantStatusClaimed
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUpdateEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UpdateEventRequest{
				Id: ev.ID,
				Input: &eventv1.EventInput{
					Title:       "shinnenkai",
					TotalAmount: 50000,
					TierCount:   1,
					HeldAt:      timestamppb.New(time.Now().Add(24 * time.Hour)),
					Tiers:       []*eventv1.TierConfig{{Tier: 1, Count: 5}},
				},
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusBadRequest)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeFailedPrecondition, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorLocked()), localized)
	})
}

func TestEventService_DeleteEvent(t *testing.T) {
	t.Parallel()

	t.Run("deletes event successfully", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		ev := fixture.Event(func(m *model.Event) { m.UserID = userID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceDeleteEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.DeleteEventRequest{Id: ev.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK)
	})

	t.Run("returns not found for nonexistent event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceDeleteEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.DeleteEventRequest{Id: "nonexistent"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusNotFound)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeNotFound, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorNotFound()), localized)
	})

	t.Run("returns forbidden for other user's event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)
		otherUser := tseed.EndUser(t, infra.WriterDB)
		ev := fixture.Event(func(m *model.Event) { m.UserID = otherUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceDeleteEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.DeleteEventRequest{Id: ev.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusForbidden)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodePermissionDenied, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorForbidden()), localized)
	})
}

func TestEventService_ListEventParticipants(t *testing.T) {
	t.Parallel()

	t.Run("returns participants for owned event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		ev := fixture.Event(func(m *model.Event) { m.UserID = userID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		p := fixture.EventParticipant(func(m *model.EventParticipant) {
			m.EventID = ev.ID
			m.Amount = 5000
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		// act
		var out eventv1.ListEventParticipantsResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceListEventParticipantsProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.ListEventParticipantsRequest{EventId: ev.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		require.Len(t, out.GetParticipants(), 1)
		assert.Equal(t, p.ID, out.GetParticipants()[0].GetId())
		assert.Equal(t, int32(5000), out.GetParticipants()[0].GetAmount())
	})

	t.Run("returns forbidden for other user's event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)
		otherUser := tseed.EndUser(t, infra.WriterDB)
		ev := fixture.Event(func(m *model.Event) { m.UserID = otherUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceListEventParticipantsProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.ListEventParticipantsRequest{EventId: ev.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusForbidden)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodePermissionDenied, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorForbidden()), localized)
	})
}

func TestEventService_UpdateParticipantStatus(t *testing.T) {
	t.Parallel()

	t.Run("confirms participant status", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		ev := fixture.Event(func(m *model.Event) { m.UserID = userID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		p := fixture.EventParticipant(func(m *model.EventParticipant) {
			m.EventID = ev.ID
			m.Status = model.ParticipantStatusClaimed
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		// act
		var out eventv1.UpdateParticipantStatusResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUpdateParticipantStatusProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UpdateParticipantStatusRequest{
				EventId:       ev.ID,
				ParticipantId: p.ID,
				Status:        eventv1.ParticipantStatus_PARTICIPANT_STATUS_CONFIRMED,
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Equal(t, eventv1.ParticipantStatus_PARTICIPANT_STATUS_CONFIRMED, out.GetParticipant().GetStatus())
	})

	t.Run("returns not found for nonexistent participant", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUpdateParticipantStatusProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UpdateParticipantStatusRequest{
				EventId:       "nonexistent",
				ParticipantId: "nonexistent",
				Status:        eventv1.ParticipantStatus_PARTICIPANT_STATUS_CONFIRMED,
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusNotFound)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeNotFound, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorParticipantNotFound()), localized)
	})

	t.Run("returns forbidden for other user's event participant", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)
		otherUser := tseed.EndUser(t, infra.WriterDB)
		ev := fixture.Event(func(m *model.Event) { m.UserID = otherUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		p := fixture.EventParticipant(func(m *model.EventParticipant) {
			m.EventID = ev.ID
			m.Status = model.ParticipantStatusClaimed
		})
		require.NoError(t, query.EventParticipants(infra.WriterDB).Create(t.Context(), &p))

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUpdateParticipantStatusProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UpdateParticipantStatusRequest{
				EventId:       ev.ID,
				ParticipantId: p.ID,
				Status:        eventv1.ParticipantStatus_PARTICIPANT_STATUS_CONFIRMED,
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusForbidden)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodePermissionDenied, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorForbidden()), localized)
	})

	t.Run("returns error for unspecified status", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUpdateParticipantStatusProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UpdateParticipantStatusRequest{
				EventId:       "some-event",
				ParticipantId: "some-participant",
				Status:        eventv1.ParticipantStatus_PARTICIPANT_STATUS_UNSPECIFIED,
			}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusBadRequest)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeInvalidArgument, connErr)
		fv := ctest.FieldViolation(t, connErr)
		assert.Equal(t, "status", fv.GetField())
	})
}

func TestEventService_ListMyEvents_Filter(t *testing.T) {
	t.Parallel()

	t.Run("returns only active events by default", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		active := fixture.Event(func(m *model.Event) { m.UserID = userID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &active))
		now := time.Now()
		archived := fixture.Event(func(m *model.Event) { m.UserID = userID; m.ArchivedAt = &now })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &archived))

		// act
		var out eventv1.ListMyEventsResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceListMyEventsProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.ListMyEventsRequest{IncludeArchived: false}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		require.Len(t, out.GetEvents(), 1)
		assert.Equal(t, active.ID, out.GetEvents()[0].GetId())
	})

	t.Run("returns only archived events when include_archived is true", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		active := fixture.Event(func(m *model.Event) { m.UserID = userID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &active))
		now := time.Now()
		archived := fixture.Event(func(m *model.Event) { m.UserID = userID; m.ArchivedAt = &now })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &archived))

		// act
		var out eventv1.ListMyEventsResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceListMyEventsProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.ListMyEventsRequest{IncludeArchived: true}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		require.Len(t, out.GetEvents(), 1)
		assert.Equal(t, archived.ID, out.GetEvents()[0].GetId())
	})
}

func TestEventService_ArchiveEvent(t *testing.T) {
	t.Parallel()

	t.Run("archives event successfully", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		ev := fixture.Event(func(m *model.Event) { m.UserID = userID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		// act
		var out eventv1.ArchiveEventResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceArchiveEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.ArchiveEventRequest{Id: ev.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Equal(t, ev.ID, out.GetEvent().GetId())
		assert.NotNil(t, out.GetEvent().GetArchivedAt())
	})

	t.Run("returns not found for nonexistent event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceArchiveEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.ArchiveEventRequest{Id: "nonexistent"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusNotFound)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeNotFound, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorNotFound()), localized)
	})

	t.Run("returns forbidden for other user's event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)
		otherUser := tseed.EndUser(t, infra.WriterDB)
		ev := fixture.Event(func(m *model.Event) { m.UserID = otherUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceArchiveEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.ArchiveEventRequest{Id: ev.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusForbidden)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodePermissionDenied, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorForbidden()), localized)
	})
}

func TestEventService_UnarchiveEvent(t *testing.T) {
	t.Parallel()

	t.Run("unarchives event successfully", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		userID, authHeader := ctest.UserSession(t, infra)
		now := time.Now()
		ev := fixture.Event(func(m *model.Event) { m.UserID = userID; m.ArchivedAt = &now })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		// act
		var out eventv1.UnarchiveEventResponse
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUnarchiveEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UnarchiveEventRequest{Id: ev.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusOK).Out(&out)
		assert.Equal(t, ev.ID, out.GetEvent().GetId())
		assert.Nil(t, out.GetEvent().GetArchivedAt())
	})

	t.Run("returns not found for nonexistent event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUnarchiveEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UnarchiveEventRequest{Id: "nonexistent"}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusNotFound)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodeNotFound, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorNotFound()), localized)
	})

	t.Run("returns forbidden for other user's event", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		_, authHeader := ctest.UserSession(t, infra)
		otherUser := tseed.EndUser(t, infra.WriterDB)
		now := time.Now()
		ev := fixture.Event(func(m *model.Event) { m.UserID = otherUser.UserID; m.ArchivedAt = &now })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		// act
		ct := contest.NewWith(t,
			contest.Bind(eventv1connect.NewEventServiceHandler)(handler.NewEventService(infra)),
			connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
		).
			Procedure(eventv1connect.EventServiceUnarchiveEventProcedure).
			Header("Authorization", authHeader).
			In(&eventv1.UnarchiveEventRequest{Id: ev.ID}).
			Do()

		// assert
		ct.ExpectStatus(http.StatusForbidden)
		connErr := ct.Err()
		ctest.AssertCode(t, connect.CodePermissionDenied, connErr)
		localized := ctest.LocalizedMessage(t, connErr)
		assert.Equal(t, i18n.Japanese(messages.EventUseCaseErrorForbidden()), localized)
	})
}
