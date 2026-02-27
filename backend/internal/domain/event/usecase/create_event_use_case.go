package usecase

import (
	"context"
	"time"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/lib/slicex"
	"github.com/mickamy/sampay/internal/lib/ulid"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

type CreateEventInput struct {
	Title       string
	Description string
	TotalAmount int
	TierCount   int
	HeldAt      time.Time
	Tiers       []TierConfig
}

type CreateEventOutput struct {
	Event model.Event
}

type CreateEvent interface {
	Do(ctx context.Context, input CreateEventInput) (CreateEventOutput, error)
}

type createEvent struct {
	_         CreateEvent          `inject:"returns"`
	_         *di.Infra            `inject:"param"`
	writer    *database.Writer     `inject:""`
	eventRepo repository.Event     `inject:""`
	tierRepo  repository.EventTier `inject:""`
}

func (uc *createEvent) Do(ctx context.Context, input CreateEventInput) (CreateEventOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	if err := validateEventInput(ctx, input.Title, input.TotalAmount, input.TierCount, input.Tiers); err != nil {
		return CreateEventOutput{}, err
	}

	tiers := make([]model.EventTier, len(input.Tiers))
	eventID := ulid.New()
	for i, tc := range input.Tiers {
		tiers[i] = model.EventTier{
			ID:      ulid.New(),
			EventID: eventID,
			Tier:    tc.Tier,
			Count:   tc.Count,
		}
	}

	ev := model.Event{
		ID:          eventID,
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		TotalAmount: input.TotalAmount,
		TierCount:   input.TierCount,
		HeldAt:      input.HeldAt,
		Tiers:       tiers,
	}
	ev.CalcTierAmounts()

	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		if err := uc.eventRepo.WithTx(tx).Create(ctx, &ev); err != nil {
			return errx.Wrap(err, "message", "failed to create event").
				WithCode(errx.Internal)
		}
		if err := uc.tierRepo.WithTx(tx).CreateAll(ctx, slicex.MapToPointer(tiers)); err != nil {
			return errx.Wrap(err, "message", "failed to create event tiers").
				WithCode(errx.Internal)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return CreateEventOutput{}, err
	}

	return CreateEventOutput{Event: ev}, nil
}
