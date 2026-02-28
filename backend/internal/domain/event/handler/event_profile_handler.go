package handler

import (
	"context"

	"connectrpc.com/connect"

	"github.com/mickamy/sampay/config"
	eventv1 "github.com/mickamy/sampay/gen/event/v1"
	"github.com/mickamy/sampay/gen/event/v1/eventv1connect"
	userv1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/event/mapper"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/usecase"
	umapper "github.com/mickamy/sampay/internal/domain/user/mapper"
	umodel "github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/lib/converter"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/lib/slicex"
)

var _ eventv1connect.EventProfileServiceHandler = (*EventProfile)(nil)

type EventProfile struct {
	_            *di.Infra            `inject:"param"`
	getEvent     usecase.GetEvent     `inject:""`
	joinEvent    usecase.JoinEvent    `inject:""`
	claimPayment usecase.ClaimPayment `inject:""`
}

func (h *EventProfile) GetEvent(
	ctx context.Context, r *connect.Request[eventv1.GetEventRequest],
) (*connect.Response[eventv1.GetEventResponse], error) {
	out, err := h.getEvent.Do(ctx, usecase.GetEventInput{
		ID: r.Msg.GetId(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	cloudfrontURL := config.AWS().CloudfrontURL()
	ev := mapper.ToV1Event(out.Event)
	user := umapper.ToV1User(out.User)
	methods := slicex.Map(out.PaymentMethods, func(m umodel.UserPaymentMethod) *userv1.PaymentMethod {
		return umapper.ToV1PaymentMethod(m, cloudfrontURL)
	})
	participants := slicex.Map(out.Event.Participants, func(p model.EventParticipant) *eventv1.EventParticipant {
		ep := mapper.ToV1EventParticipant(p)
		return &ep
	})

	return connect.NewResponse(&eventv1.GetEventResponse{
		Event:          &ev,
		User:           &user,
		PaymentMethods: methods,
		Participants:   participants,
	}), nil
}

func (h *EventProfile) JoinEvent(
	ctx context.Context, r *connect.Request[eventv1.JoinEventRequest],
) (*connect.Response[eventv1.JoinEventResponse], error) {
	out, err := h.joinEvent.Do(ctx, usecase.JoinEventInput{
		EventID: r.Msg.GetEventId(),
		Name:    r.Msg.GetName(),
		Tier:    converter.Int32ToInt(r.Msg.GetTier()),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	participant := mapper.ToV1EventParticipant(out.Participant)
	return connect.NewResponse(&eventv1.JoinEventResponse{
		Participant: &participant,
	}), nil
}

func (h *EventProfile) ClaimPayment(
	ctx context.Context, r *connect.Request[eventv1.ClaimPaymentRequest],
) (*connect.Response[eventv1.ClaimPaymentResponse], error) {
	out, err := h.claimPayment.Do(ctx, usecase.ClaimPaymentInput{
		ParticipantID: r.Msg.GetParticipantId(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	participant := mapper.ToV1EventParticipant(out.Participant)
	return connect.NewResponse(&eventv1.ClaimPaymentResponse{
		Participant: &participant,
	}), nil
}
