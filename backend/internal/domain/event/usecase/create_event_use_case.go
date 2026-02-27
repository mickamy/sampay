package usecase

import (
	"context"
	"time"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/lib/ulid"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

type CreateEventInput struct {
	Title       string
	Description string
	TotalAmount int
	TierCount   int
	HeldAt      time.Time
}

type CreateEventOutput struct {
	Event model.Event
}

type CreateEvent interface {
	Do(ctx context.Context, input CreateEventInput) (CreateEventOutput, error)
}

type createEvent struct {
	_         CreateEvent      `inject:"returns"`
	_         *di.Infra        `inject:"param"`
	writer    *database.Writer `inject:""`
	eventRepo repository.Event `inject:""`
}

func (uc *createEvent) Do(ctx context.Context, input CreateEventInput) (CreateEventOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	if err := validateEventInput(ctx, input.Title, input.TotalAmount, input.TierCount); err != nil {
		return CreateEventOutput{}, err
	}

	ev := model.Event{
		ID:          ulid.New(),
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		TotalAmount: input.TotalAmount,
		TierCount:   input.TierCount,
		HeldAt:      input.HeldAt,
	}

	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		if err := uc.eventRepo.WithTx(tx).Create(ctx, &ev); err != nil {
			return errx.Wrap(err, "message", "failed to create event", "event", ev).
				WithCode(errx.Internal)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return CreateEventOutput{}, err
	}

	return CreateEventOutput{Event: ev}, nil
}
