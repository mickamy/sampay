package usecase

import (
	"context"
	"time"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/lib/ulid"
	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var (
	ErrCreateEventEmptyTitle = cmodel.NewLocalizableError(
		errx.NewSentinel("title is required", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorTitleRequired())
	ErrCreateEventNegativeTotalAmount = cmodel.NewLocalizableError(
		errx.NewSentinel("total_amount must be positive", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorTotalAmountPositive())
	ErrCreateEventInvalidTierCount = cmodel.NewLocalizableError(
		errx.NewSentinel("tier_count must be 1, 3, or 5", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorTierCountInvalid())
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

	if err := uc.validateEventInput(ctx, input.Title, input.TotalAmount, input.TierCount); err != nil {
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

func (uc *createEvent) validateEventInput(ctx context.Context, title string, totalAmount, tierCount int) error {
	if title == "" {
		return errx.Wrap(ErrCreateEventEmptyTitle).
			WithFieldViolation("title", ErrCreateEventEmptyTitle.LocalizeContext(ctx))
	}
	if totalAmount <= 0 {
		return errx.Wrap(ErrCreateEventNegativeTotalAmount, "total_amount", totalAmount).
			WithFieldViolation("total_amount", ErrCreateEventNegativeTotalAmount.LocalizeContext(ctx))
	}
	switch tierCount {
	case 1, 3, 5:
		// valid
	default:
		return errx.Wrap(ErrCreateEventInvalidTierCount, "tier_count", tierCount).
			WithFieldViolation("tier_count", ErrCreateEventInvalidTierCount.LocalizeContext(ctx))
	}
	return nil
}
