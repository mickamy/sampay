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
	"github.com/mickamy/sampay/internal/lib/slicex"
	"github.com/mickamy/sampay/internal/lib/ulid"
	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var (
	ErrUpdateEventNotFound = cmodel.NewLocalizableError(
		errx.NewSentinel("event not found", errx.NotFound),
	).WithMessages(messages.EventUseCaseErrorNotFound())
	ErrUpdateEventForbidden = cmodel.NewLocalizableError(
		errx.NewSentinel("forbidden", errx.PermissionDenied),
	).WithMessages(messages.EventUseCaseErrorForbidden())
	ErrUpdateEventLocked = cmodel.NewLocalizableError(
		errx.NewSentinel("event is locked", errx.FailedPrecondition),
	).WithMessages(messages.EventUseCaseErrorLocked())
	ErrUpdateEventInvalidTierUpdate = cmodel.NewLocalizableError(
		errx.NewSentinel("invalid tier update", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorInvalidTierUpdate())
)

type UpdateEventInput struct {
	ID          string
	Title       string
	Description string
	TotalAmount int
	TierCount   int
	HeldAt      time.Time
	Tiers       []TierConfig
}

type UpdateEventOutput struct {
	Event model.Event
}

type UpdateEvent interface {
	Do(ctx context.Context, input UpdateEventInput) (UpdateEventOutput, error)
}

type updateEvent struct {
	_               UpdateEvent                 `inject:"returns"`
	_               *di.Infra                   `inject:"param"`
	writer          *database.Writer            `inject:""`
	eventRepo       repository.Event            `inject:""`
	tierRepo        repository.EventTier        `inject:""`
	participantRepo repository.EventParticipant `inject:""`
}

func (uc *updateEvent) Do(ctx context.Context, input UpdateEventInput) (UpdateEventOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	if err := validateEventInput(
		ctx, input.Title, input.TotalAmount, input.TierCount, input.HeldAt, input.Tiers,
	); err != nil {
		return UpdateEventOutput{}, err
	}

	var ev model.Event
	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		var err error
		ev, err = uc.eventRepo.WithTx(tx).Get(
			ctx, input.ID, repository.EventPreloadParticipants(),
		)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ErrUpdateEventNotFound
			}
			return errx.Wrap(err, "message", "failed to get event", "id", input.ID).
				WithCode(errx.Internal)
		}

		if ev.UserID != userID {
			return ErrUpdateEventForbidden
		}

		for _, p := range ev.Participants {
			if p.Status != model.ParticipantStatusUnpaid {
				return ErrUpdateEventLocked
			}
		}

		tiers := make([]model.EventTier, len(input.Tiers))
		for i, tc := range input.Tiers {
			tiers[i] = model.EventTier{
				ID:      ulid.New(),
				EventID: ev.ID,
				Tier:    tc.Tier,
				Count:   tc.Count,
			}
		}

		// Build a set of new tier numbers for validation.
		newTierSet := make(map[int]bool, len(input.Tiers))
		for _, tc := range input.Tiers {
			newTierSet[tc.Tier] = true
		}
		for _, p := range ev.Participants {
			if !newTierSet[p.Tier] {
				return ErrUpdateEventInvalidTierUpdate
			}
		}

		ev.Title = input.Title
		ev.Description = input.Description
		ev.TotalAmount = input.TotalAmount
		ev.TierCount = input.TierCount
		ev.HeldAt = input.HeldAt
		ev.Tiers = tiers
		ev.CalcTierAmounts()

		if err := uc.eventRepo.WithTx(tx).Update(ctx, &ev); err != nil {
			return errx.Wrap(err, "message", "failed to update event", "id", input.ID).
				WithCode(errx.Internal)
		}

		if err := uc.tierRepo.WithTx(tx).DeleteByEventID(ctx, ev.ID); err != nil {
			return errx.Wrap(err, "message", "failed to delete old tiers", "id", ev.ID).
				WithCode(errx.Internal)
		}

		if err := uc.tierRepo.WithTx(tx).CreateAll(ctx, slicex.MapToPointer(tiers)); err != nil {
			return errx.Wrap(err, "message", "failed to create tiers", "id", ev.ID).
				WithCode(errx.Internal)
		}

		for i := range ev.Participants {
			ev.Participants[i].Amount = ev.TierAmount(ev.Participants[i].Tier)
			if err := uc.participantRepo.WithTx(tx).Update(ctx, &ev.Participants[i]); err != nil {
				return errx.Wrap(err, "message", "failed to update participant amount").
					WithCode(errx.Internal)
			}
		}

		ev.Tiers = tiers
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return UpdateEventOutput{}, err
	}

	return UpdateEventOutput{Event: ev}, nil
}
