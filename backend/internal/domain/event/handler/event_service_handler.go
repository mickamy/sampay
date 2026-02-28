package handler

import (
	"context"
	"time"

	"connectrpc.com/connect"
	"github.com/mickamy/errx"

	v1 "github.com/mickamy/sampay/gen/event/v1"
	"github.com/mickamy/sampay/gen/event/v1/eventv1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/event/mapper"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/usecase"
	"github.com/mickamy/sampay/internal/lib/converter"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/lib/slicex"
)

var _ eventv1connect.EventServiceHandler = (*EventService)(nil)

type EventService struct {
	_                       *di.Infra                       `inject:"param"`
	listMyEvents            usecase.ListMyEvents            `inject:""`
	createEvent             usecase.CreateEvent             `inject:""`
	updateEvent             usecase.UpdateEvent             `inject:""`
	deleteEvent             usecase.DeleteEvent             `inject:""`
	listEventParticipants   usecase.ListEventParticipants   `inject:""`
	updateParticipantStatus usecase.UpdateParticipantStatus `inject:""`
	archiveEvent            usecase.ArchiveEvent            `inject:""`
	unarchiveEvent          usecase.UnarchiveEvent          `inject:""`
}

func (h *EventService) ListMyEvents(
	ctx context.Context, r *connect.Request[v1.ListMyEventsRequest],
) (*connect.Response[v1.ListMyEventsResponse], error) {
	out, err := h.listMyEvents.Do(ctx, usecase.ListMyEventsInput{
		IncludeArchived: r.Msg.GetIncludeArchived(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	events := slicex.Map(out.Events, func(e model.Event) *v1.Event {
		ev := mapper.ToV1Event(e)
		return &ev
	})
	return connect.NewResponse(&v1.ListMyEventsResponse{
		Events: events,
	}), nil
}

func (h *EventService) CreateEvent(
	ctx context.Context, r *connect.Request[v1.CreateEventRequest],
) (*connect.Response[v1.CreateEventResponse], error) {
	input := r.Msg.GetInput()
	tiers := slicex.Map(input.GetTiers(), func(tc *v1.TierConfig) usecase.TierConfig {
		return usecase.TierConfig{
			Tier:  converter.Int32ToInt(tc.GetTier()),
			Count: converter.Int32ToInt(tc.GetCount()),
		}
	})

	var heldAt time.Time
	if ts := input.GetHeldAt(); ts != nil {
		heldAt = ts.AsTime()
	}

	out, err := h.createEvent.Do(ctx, usecase.CreateEventInput{
		Title:       input.GetTitle(),
		Description: input.GetDescription(),
		TotalAmount: converter.Int32ToInt(input.GetTotalAmount()),
		TierCount:   converter.Int32ToInt(input.GetTierCount()),
		HeldAt:      heldAt,
		Tiers:       tiers,
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	ev := mapper.ToV1Event(out.Event)
	return connect.NewResponse(&v1.CreateEventResponse{
		Event: &ev,
	}), nil
}

func (h *EventService) UpdateEvent(
	ctx context.Context, r *connect.Request[v1.UpdateEventRequest],
) (*connect.Response[v1.UpdateEventResponse], error) {
	input := r.Msg.GetInput()
	tiers := slicex.Map(input.GetTiers(), func(tc *v1.TierConfig) usecase.TierConfig {
		return usecase.TierConfig{
			Tier:  converter.Int32ToInt(tc.GetTier()),
			Count: converter.Int32ToInt(tc.GetCount()),
		}
	})

	var updateHeldAt time.Time
	if ts := input.GetHeldAt(); ts != nil {
		updateHeldAt = ts.AsTime()
	}

	out, err := h.updateEvent.Do(ctx, usecase.UpdateEventInput{
		ID:          r.Msg.GetId(),
		Title:       input.GetTitle(),
		Description: input.GetDescription(),
		TotalAmount: converter.Int32ToInt(input.GetTotalAmount()),
		TierCount:   converter.Int32ToInt(input.GetTierCount()),
		HeldAt:      updateHeldAt,
		Tiers:       tiers,
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	ev := mapper.ToV1Event(out.Event)
	return connect.NewResponse(&v1.UpdateEventResponse{
		Event: &ev,
	}), nil
}

func (h *EventService) DeleteEvent(
	ctx context.Context, r *connect.Request[v1.DeleteEventRequest],
) (*connect.Response[v1.DeleteEventResponse], error) {
	_, err := h.deleteEvent.Do(ctx, usecase.DeleteEventInput{
		ID: r.Msg.GetId(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	return connect.NewResponse(&v1.DeleteEventResponse{}), nil
}

func (h *EventService) ListEventParticipants(
	ctx context.Context, r *connect.Request[v1.ListEventParticipantsRequest],
) (*connect.Response[v1.ListEventParticipantsResponse], error) {
	out, err := h.listEventParticipants.Do(ctx, usecase.ListEventParticipantsInput{
		EventID: r.Msg.GetEventId(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	participants := slicex.Map(out.Participants, func(p model.EventParticipant) *v1.EventParticipant {
		ep := mapper.ToV1EventParticipant(p)
		return &ep
	})
	return connect.NewResponse(&v1.ListEventParticipantsResponse{
		Participants: participants,
	}), nil
}

func (h *EventService) UpdateParticipantStatus(
	ctx context.Context, r *connect.Request[v1.UpdateParticipantStatusRequest],
) (*connect.Response[v1.UpdateParticipantStatusResponse], error) {
	status, err := converter.FromV1ParticipantStatus(r.Msg.GetStatus())
	if err != nil {
		return nil, errx.Wrap(err, "message", "invalid status").
			WithCode(errx.InvalidArgument).
			WithFieldViolation("status", err.Error())
	}

	out, err := h.updateParticipantStatus.Do(ctx, usecase.UpdateParticipantStatusInput{
		EventID:       r.Msg.GetEventId(),
		ParticipantID: r.Msg.GetParticipantId(),
		Status:        status,
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	participant := mapper.ToV1EventParticipant(out.Participant)
	return connect.NewResponse(&v1.UpdateParticipantStatusResponse{
		Participant: &participant,
	}), nil
}

func (h *EventService) ArchiveEvent(
	ctx context.Context, r *connect.Request[v1.ArchiveEventRequest],
) (*connect.Response[v1.ArchiveEventResponse], error) {
	out, err := h.archiveEvent.Do(ctx, usecase.ArchiveEventInput{
		ID: r.Msg.GetId(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	ev := mapper.ToV1Event(out.Event)
	return connect.NewResponse(&v1.ArchiveEventResponse{
		Event: &ev,
	}), nil
}

func (h *EventService) UnarchiveEvent(
	ctx context.Context, r *connect.Request[v1.UnarchiveEventRequest],
) (*connect.Response[v1.UnarchiveEventResponse], error) {
	out, err := h.unarchiveEvent.Do(ctx, usecase.UnarchiveEventInput{
		ID: r.Msg.GetId(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	ev := mapper.ToV1Event(out.Event)
	return connect.NewResponse(&v1.UnarchiveEventResponse{
		Event: &ev,
	}), nil
}
