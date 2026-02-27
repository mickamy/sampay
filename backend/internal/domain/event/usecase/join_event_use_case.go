package usecase

import (
	"context"
	"errors"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/lib/ulid"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var (
	ErrJoinEventNotFound = cmodel.NewLocalizableError(
		errx.NewSentinel("event not found", errx.NotFound),
	).WithMessages(messages.EventUseCaseErrorNotFound())
	ErrJoinEventEmptyName = cmodel.NewLocalizableError(
		errx.NewSentinel("name is required", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorNameRequired())
	ErrJoinEventInvalidTier = cmodel.NewLocalizableError(
		errx.NewSentinel("invalid tier", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorInvalidTier())
)

type JoinEventInput struct {
	EventID string
	Name    string
	Tier    int
}

type JoinEventOutput struct {
	Participant model.EventParticipant
}

type JoinEvent interface {
	Do(ctx context.Context, input JoinEventInput) (JoinEventOutput, error)
}

type joinEvent struct {
	_               JoinEvent                   `inject:"returns"`
	_               *di.Infra                   `inject:"param"`
	writer          *database.Writer            `inject:""`
	eventRepo       repository.Event            `inject:""`
	participantRepo repository.EventParticipant `inject:""`
}

func (uc *joinEvent) Do(ctx context.Context, input JoinEventInput) (JoinEventOutput, error) {
	if input.Name == "" {
		return JoinEventOutput{}, ErrJoinEventEmptyName
	}

	var participant model.EventParticipant

	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		ev, err := uc.eventRepo.WithTx(tx).Get(ctx, input.EventID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ErrJoinEventNotFound
			}
			return errx.Wrap(err, "message", "failed to get event", "event_id", input.EventID).
				WithCode(errx.Internal)
		}

		if input.Tier < 1 || input.Tier > ev.TierCount {
			return ErrJoinEventInvalidTier
		}

		participant = model.EventParticipant{
			ID:      ulid.New(),
			EventID: input.EventID,
			Name:    input.Name,
			Tier:    input.Tier,
			Status:  model.ParticipantStatusUnpaid,
		}

		if err := uc.participantRepo.WithTx(tx).Create(ctx, &participant); err != nil {
			return errx.Wrap(err, "message", "failed to create participant").
				WithCode(errx.Internal)
		}

		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return JoinEventOutput{}, err
	}

	return JoinEventOutput{Participant: participant}, nil
}
