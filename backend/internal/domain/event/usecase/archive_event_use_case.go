package usecase

import (
	"context"
	"errors"
	"time"

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
	ErrArchiveEventNotFound = cmodel.NewLocalizableError(
		errx.NewSentinel("event not found", errx.NotFound),
	).WithMessages(messages.EventUseCaseErrorNotFound())
	ErrArchiveEventForbidden = cmodel.NewLocalizableError(
		errx.NewSentinel("forbidden", errx.PermissionDenied),
	).WithMessages(messages.EventUseCaseErrorForbidden())
)

type ArchiveEventInput struct {
	ID string
}

type ArchiveEventOutput struct {
	Event model.Event
}

type ArchiveEvent interface {
	Do(ctx context.Context, input ArchiveEventInput) (ArchiveEventOutput, error)
}

type archiveEvent struct {
	_         ArchiveEvent     `inject:"returns"`
	_         *di.Infra        `inject:"param"`
	writer    *database.Writer `inject:""`
	eventRepo repository.Event `inject:""`
}

func (uc *archiveEvent) Do(ctx context.Context, input ArchiveEventInput) (ArchiveEventOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	var ev model.Event
	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		var err error
		ev, err = uc.eventRepo.WithTx(tx).Get(ctx, input.ID)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ErrArchiveEventNotFound
			}
			return errx.Wrap(err, "message", "failed to get event", "id", input.ID).
				WithCode(errx.Internal)
		}

		if ev.UserID != userID {
			return ErrArchiveEventForbidden
		}

		now := time.Now()
		ev.ArchivedAt = &now
		if err := uc.eventRepo.WithTx(tx).Update(ctx, &ev); err != nil {
			return errx.Wrap(err, "message", "failed to archive event", "id", input.ID).
				WithCode(errx.Internal)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return ArchiveEventOutput{}, err
	}

	return ArchiveEventOutput{Event: ev}, nil
}
