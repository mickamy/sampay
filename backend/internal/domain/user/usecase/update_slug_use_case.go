package usecase

import (
	"context"
	"errors"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

type UpdateSlugInput struct {
	Slug string
}

type UpdateSlugOutput struct {
	User model.EndUser
}

type UpdateSlug interface {
	Do(ctx context.Context, input UpdateSlugInput) (UpdateSlugOutput, error)
}

type updateSlug struct {
	_           UpdateSlug         `inject:"returns"`
	_           *di.Infra          `inject:"param"`
	writer      *database.Writer   `inject:""`
	endUserRepo repository.EndUser `inject:""`
}

func (uc *updateSlug) Do(
	ctx context.Context,
	input UpdateSlugInput,
) (UpdateSlugOutput, error) {
	if err := model.ValidateSlug(input.Slug); err != nil {
		return UpdateSlugOutput{}, errx.Wrap(err, "message", "slug validation failed").
			WithCode(errx.InvalidArgument)
	}

	userID := contexts.MustAuthenticatedUserID(ctx)

	var endUser model.EndUser
	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		// Check uniqueness.
		_, err := uc.endUserRepo.WithTx(tx).GetBySlug(ctx, input.Slug)
		if err == nil {
			return errx.Wrap(model.ErrSlugAlreadyTaken, "message", "slug already taken", "slug", input.Slug).
				WithCode(errx.AlreadyExists)
		}
		if !errors.Is(err, database.ErrNotFound) {
			return errx.Wrap(err, "message", "failed to check slug uniqueness", "slug", input.Slug).
				WithCode(errx.Internal)
		}

		endUser, err = uc.endUserRepo.WithTx(tx).Get(ctx, userID)
		if err != nil {
			return errx.Wrap(err, "message", "failed to get end user", "user_id", userID).
				WithCode(errx.Internal)
		}

		endUser.Slug = input.Slug
		if err := uc.endUserRepo.WithTx(tx).Update(ctx, &endUser); err != nil {
			return errx.Wrap(err, "message", "failed to update slug", "user_id", userID).
				WithCode(errx.Internal)
		}

		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return UpdateSlugOutput{}, err
	}

	return UpdateSlugOutput{User: endUser}, nil
}
