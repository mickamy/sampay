package usecase

import (
	"context"
	"errors"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/domain/event/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var (
	ErrDeleteEventNotFound = cmodel.NewLocalizableError(
		errx.NewSentinel("event not found", errx.NotFound),
	).WithMessages(messages.EventUseCaseErrorNotFound())
	ErrDeleteEventForbidden = cmodel.NewLocalizableError(
		errx.NewSentinel("forbidden", errx.PermissionDenied),
	).WithMessages(messages.EventUseCaseErrorForbidden())
)

type DeleteEventInput struct {
	ID string
}

type DeleteEventOutput struct{}

type DeleteEvent interface {
	Do(ctx context.Context, input DeleteEventInput) (DeleteEventOutput, error)
}

type deleteEvent struct {
	_         DeleteEvent      `inject:"returns"`
	_         *di.Infra        `inject:"param"`
	writer    *database.Writer `inject:""`
	eventRepo repository.Event `inject:""`
}

func (uc *deleteEvent) Do(ctx context.Context, input DeleteEventInput) (DeleteEventOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		ev, err := uc.eventRepo.WithTx(tx).Get(ctx, input.ID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ErrDeleteEventNotFound
			}
			return errx.Wrap(err, "message", "failed to get event", "id", input.ID).
				WithCode(errx.Internal)
		}

		if ev.UserID != userID {
			return ErrDeleteEventForbidden
		}

		if err := uc.eventRepo.WithTx(tx).Delete(ctx, input.ID); err != nil {
			return errx.Wrap(err, "message", "failed to delete event", "id", input.ID).
				WithCode(errx.Internal)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return DeleteEventOutput{}, err
	}

	return DeleteEventOutput{}, nil
}
