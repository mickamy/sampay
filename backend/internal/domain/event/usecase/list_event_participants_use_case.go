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
	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var (
	ErrListEventParticipantsNotFound = cmodel.NewLocalizableError(
		errx.NewSentinel("event not found", errx.NotFound),
	).WithMessages(messages.EventUseCaseErrorNotFound())
	ErrListEventParticipantsForbidden = cmodel.NewLocalizableError(
		errx.NewSentinel("forbidden", errx.PermissionDenied),
	).WithMessages(messages.EventUseCaseErrorForbidden())
)

type ListEventParticipantsInput struct {
	EventID string
}

type ListEventParticipantsOutput struct {
	Participants []model.EventParticipant
}

type ListEventParticipants interface {
	Do(ctx context.Context, input ListEventParticipantsInput) (ListEventParticipantsOutput, error)
}

type listEventParticipants struct {
	_               ListEventParticipants      `inject:"returns"`
	_               *di.Infra                  `inject:"param"`
	reader          *database.Reader           `inject:""`
	eventRepo       repository.Event           `inject:""`
	participantRepo repository.EventParticipant `inject:""`
}

func (uc *listEventParticipants) Do(
	ctx context.Context, input ListEventParticipantsInput,
) (ListEventParticipantsOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	var participants []model.EventParticipant

	if err := uc.reader.Transaction(ctx, func(tx *database.DB) error {
		ev, err := uc.eventRepo.WithTx(tx).Get(ctx, input.EventID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ErrListEventParticipantsNotFound
			}
			return errx.Wrap(err, "message", "failed to get event", "id", input.EventID).
				WithCode(errx.Internal)
		}

		if ev.UserID != userID {
			return ErrListEventParticipantsForbidden
		}

		participants, err = uc.participantRepo.WithTx(tx).ListByEventID(ctx, input.EventID)
		if err != nil {
			return errx.Wrap(err, "message", "failed to list participants", "event_id", input.EventID).
				WithCode(errx.Internal)
		}

		model.CalcAmounts(ev.TotalAmount, participants)

		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return ListEventParticipantsOutput{}, err
	}

	return ListEventParticipantsOutput{Participants: participants}, nil
}
