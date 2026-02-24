package usecase

import (
	"context"
	"errors"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

type CheckSlugAvailabilityInput struct {
	Slug string
}

type CheckSlugAvailabilityOutput struct {
	Available bool
}

type CheckSlugAvailability interface {
	Do(ctx context.Context, input CheckSlugAvailabilityInput) (CheckSlugAvailabilityOutput, error)
}

type checkSlugAvailability struct {
	_           CheckSlugAvailability `inject:"returns"`
	_           *di.Infra             `inject:"param"`
	reader      *database.Reader      `inject:""`
	endUserRepo repository.EndUser    `inject:""`
}

func (uc *checkSlugAvailability) Do(
	ctx context.Context,
	input CheckSlugAvailabilityInput,
) (CheckSlugAvailabilityOutput, error) {
	if err := model.ValidateSlug(input.Slug); err != nil {
		return CheckSlugAvailabilityOutput{Available: false}, nil
	}

	var available bool
	if err := uc.reader.Transaction(ctx, func(tx *database.DB) error {
		_, err := uc.endUserRepo.WithTx(tx).GetBySlug(ctx, input.Slug)
		if errors.Is(err, database.ErrNotFound) {
			available = true
			return nil
		}
		if err != nil {
			return errx.Wrap(err, "message", "failed to check slug", "slug", input.Slug).
				WithCode(errx.Internal)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return CheckSlugAvailabilityOutput{}, err
	}

	return CheckSlugAvailabilityOutput{Available: available}, nil
}
