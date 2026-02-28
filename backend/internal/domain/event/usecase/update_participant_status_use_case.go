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
	ErrUpdateParticipantStatusNotFound = cmodel.NewLocalizableError(
		errx.NewSentinel("participant not found", errx.NotFound),
	).WithMessages(messages.EventUseCaseErrorParticipantNotFound())
	ErrUpdateParticipantStatusForbidden = cmodel.NewLocalizableError(
		errx.NewSentinel("forbidden", errx.PermissionDenied),
	).WithMessages(messages.EventUseCaseErrorForbidden())
	ErrUpdateParticipantStatusEventMismatch = cmodel.NewLocalizableError(
		errx.NewSentinel("event_id mismatch", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorEventMismatch())
)

type UpdateParticipantStatusInput struct {
	EventID       string
	ParticipantID string
	Status        model.ParticipantStatus
}

type UpdateParticipantStatusOutput struct {
	Participant model.EventParticipant
}

type UpdateParticipantStatus interface {
	Do(ctx context.Context, input UpdateParticipantStatusInput) (UpdateParticipantStatusOutput, error)
}

type updateParticipantStatus struct {
	_               UpdateParticipantStatus     `inject:"returns"`
	_               *di.Infra                   `inject:"param"`
	writer          *database.Writer            `inject:""`
	eventRepo       repository.Event            `inject:""`
	participantRepo repository.EventParticipant `inject:""`
}

func (uc *updateParticipantStatus) Do(
	ctx context.Context, input UpdateParticipantStatusInput,
) (UpdateParticipantStatusOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	var participant model.EventParticipant

	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		var err error
		participant, err = uc.participantRepo.WithTx(tx).Get(ctx, input.ParticipantID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ErrUpdateParticipantStatusNotFound
			}
			return errx.Wrap(err, "message", "failed to get participant", "id", input.ParticipantID).
				WithCode(errx.Internal)
		}

		if participant.EventID != input.EventID {
			return ErrUpdateParticipantStatusEventMismatch
		}

		ev, err := uc.eventRepo.WithTx(tx).Get(ctx, participant.EventID)
		if err != nil {
			return errx.Wrap(err, "message", "failed to get event", "id", participant.EventID).
				WithCode(errx.Internal)
		}

		if ev.UserID != userID {
			return ErrUpdateParticipantStatusForbidden
		}

		participant.Status = input.Status
		if err := uc.participantRepo.WithTx(tx).Update(ctx, &participant); err != nil {
			return errx.Wrap(err, "message", "failed to update participant status").
				WithCode(errx.Internal)
		}

		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return UpdateParticipantStatusOutput{}, err
	}

	return UpdateParticipantStatusOutput{Participant: participant}, nil
}
