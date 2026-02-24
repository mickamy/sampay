package usecase

import (
	"context"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

type GetMeInput struct{}

type GetMeOutput struct {
	User model.EndUser
}

type GetMe interface {
	Do(ctx context.Context, input GetMeInput) (GetMeOutput, error)
}

type getMe struct {
	_           GetMe              `inject:"returns"`
	_           *di.Infra          `inject:"param"`
	reader      *database.Reader   `inject:""`
	endUserRepo repository.EndUser `inject:""`
}

func (uc *getMe) Do(
	ctx context.Context,
	_ GetMeInput,
) (GetMeOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	var endUser model.EndUser
	if err := uc.reader.Transaction(ctx, func(tx *database.DB) error {
		var err error
		endUser, err = uc.endUserRepo.WithTx(tx).Get(ctx, userID)
		if err != nil {
			return errx.Wrap(err, "message", "failed to get end user", "user_id", userID).
				WithCode(errx.Internal)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return GetMeOutput{}, err
	}

	return GetMeOutput{User: endUser}, nil
}
