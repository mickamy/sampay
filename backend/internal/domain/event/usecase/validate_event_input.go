package usecase

import (
	"context"

	"github.com/mickamy/errx"

	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var (
	ErrValidateEventEmptyTitle = cmodel.NewLocalizableError(
		errx.NewSentinel("title is required", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorTitleRequired())
	ErrValidateEventNonPositiveTotalAmount = cmodel.NewLocalizableError(
		errx.NewSentinel("total_amount must be positive", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorTotalAmountPositive())
	ErrValidateEventInvalidTierCount = cmodel.NewLocalizableError(
		errx.NewSentinel("tier_count must be 1, 3, or 5", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorTierCountInvalid())
)

func validateEventInput(ctx context.Context, title string, totalAmount, tierCount int) error {
	if title == "" {
		return errx.Wrap(ErrValidateEventEmptyTitle).
			WithFieldViolation("title", ErrValidateEventEmptyTitle.LocalizeContext(ctx))
	}
	if totalAmount <= 0 {
		return errx.Wrap(ErrValidateEventNonPositiveTotalAmount, "total_amount", totalAmount).
			WithFieldViolation("total_amount", ErrValidateEventNonPositiveTotalAmount.LocalizeContext(ctx))
	}
	switch tierCount {
	case 1, 3, 5:
		// valid
	default:
		return errx.Wrap(ErrValidateEventInvalidTierCount, "tier_count", tierCount).
			WithFieldViolation("tier_count", ErrValidateEventInvalidTierCount.LocalizeContext(ctx))
	}
	return nil
}
