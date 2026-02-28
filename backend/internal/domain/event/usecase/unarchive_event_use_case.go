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
	ErrUnarchiveEventNotFound = cmodel.NewLocalizableError(
		errx.NewSentinel("event not found", errx.NotFound),
	).WithMessages(messages.EventUseCaseErrorNotFound())
	ErrUnarchiveEventForbidden = cmodel.NewLocalizableError(
		errx.NewSentinel("forbidden", errx.PermissionDenied),
	).WithMessages(messages.EventUseCaseErrorForbidden())
)

type UnarchiveEventInput struct {
	ID string
}

type UnarchiveEventOutput struct {
	Event model.Event
}

type UnarchiveEvent interface {
	Do(ctx context.Context, input UnarchiveEventInput) (UnarchiveEventOutput, error)
}

type unarchiveEvent struct {
	_         UnarchiveEvent   `inject:"returns"`
	_         *di.Infra        `inject:"param"`
	writer    *database.Writer `inject:""`
	eventRepo repository.Event `inject:""`
}

func (uc *unarchiveEvent) Do(ctx context.Context, input UnarchiveEventInput) (UnarchiveEventOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	var ev model.Event
	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		var err error
		ev, err = uc.eventRepo.WithTx(tx).Get(ctx, input.ID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ErrUnarchiveEventNotFound
			}
			return errx.Wrap(err, "message", "failed to get event", "id", input.ID).
				WithCode(errx.Internal)
		}

		if ev.UserID != userID {
			return ErrUnarchiveEventForbidden
		}

		ev.ArchivedAt = nil
		if err := uc.eventRepo.WithTx(tx).Update(ctx, &ev); err != nil {
			return errx.Wrap(err, "message", "failed to unarchive event", "id", input.ID).
				WithCode(errx.Internal)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return UnarchiveEventOutput{}, err
	}

	return UnarchiveEventOutput{Event: ev}, nil
}
