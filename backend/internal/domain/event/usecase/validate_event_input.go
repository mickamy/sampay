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
	ErrValidateEventTiersRequired = cmodel.NewLocalizableError(
		errx.NewSentinel("tiers are required", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorTiersRequired())
	ErrValidateEventInvalidTierConfig = cmodel.NewLocalizableError(
		errx.NewSentinel("invalid tier config", errx.InvalidArgument),
	).WithMessages(messages.EventUseCaseErrorInvalidTierConfig())
)

// TierConfig represents a tier configuration input (tier number + count of people).
type TierConfig struct {
	Tier  int
	Count int
}

func validateEventInput(
	ctx context.Context, title string, totalAmount, tierCount int, tiers []TierConfig,
) error {
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

	if len(tiers) != tierCount {
		return errx.Wrap(ErrValidateEventTiersRequired, "tiers", len(tiers), "tier_count", tierCount).
			WithFieldViolation("tiers", ErrValidateEventTiersRequired.LocalizeContext(ctx))
	}

	seen := make(map[int]bool, tierCount)
	for _, tc := range tiers {
		if tc.Tier < 1 || tc.Tier > tierCount || tc.Count <= 0 || seen[tc.Tier] {
			return errx.Wrap(ErrValidateEventInvalidTierConfig, "tier", tc.Tier, "count", tc.Count).
				WithFieldViolation("tiers", ErrValidateEventInvalidTierConfig.LocalizeContext(ctx))
		}
		seen[tc.Tier] = true
	}

	return nil
}
