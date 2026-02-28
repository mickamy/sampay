package usecase

import (
	"context"

	"github.com/mickamy/errx"
	"github.com/mickamy/ormgen/scope"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

type ListMyEventsInput struct {
	IncludeArchived bool
}

type ListMyEventsOutput struct {
	Events []model.Event
}

type ListMyEvents interface {
	Do(ctx context.Context, input ListMyEventsInput) (ListMyEventsOutput, error)
}

type listMyEvents struct {
	_         ListMyEvents     `inject:"returns"`
	_         *di.Infra        `inject:"param"`
	reader    *database.Reader `inject:""`
	eventRepo repository.Event `inject:""`
}

func (uc *listMyEvents) Do(ctx context.Context, input ListMyEventsInput) (ListMyEventsOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	var filterScope scope.Scope
	if input.IncludeArchived {
		filterScope = repository.EventArchivedOnly()
	} else {
		filterScope = repository.EventActiveOnly()
	}

	var events []model.Event
	if err := uc.reader.Transaction(ctx, func(tx *database.DB) error {
		var err error
		events, err = uc.eventRepo.WithTx(tx).ListByUserID(ctx, userID, filterScope)
		if err != nil {
			return errx.Wrap(err, "message", "failed to list events", "user_id", userID).
				WithCode(errx.Internal)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return ListMyEventsOutput{}, err
	}

	return ListMyEventsOutput{Events: events}, nil
}
