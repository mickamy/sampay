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

type GetUserProfileInput struct {
	Slug string
}

type GetUserProfileOutput struct {
	User           model.EndUser
	PaymentMethods []model.UserPaymentMethod
}

type GetUserProfile interface {
	Do(ctx context.Context, input GetUserProfileInput) (GetUserProfileOutput, error)
}

type getUserProfile struct {
	_                 GetUserProfile               `inject:"returns"`
	_                 *di.Infra                    `inject:"param"`
	reader            *database.Reader             `inject:""`
	endUserRepo       repository.EndUser           `inject:""`
	paymentMethodRepo repository.UserPaymentMethod `inject:""`
}

func (uc *getUserProfile) Do(ctx context.Context, input GetUserProfileInput) (GetUserProfileOutput, error) {
	var (
		endUser model.EndUser
		methods []model.UserPaymentMethod
	)

	if err := uc.reader.Transaction(ctx, func(tx *database.DB) error {
		var err error
		endUser, err = uc.endUserRepo.WithTx(tx).GetBySlug(ctx, input.Slug)
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return errx.Wrap(err, "user not found", "slug", input.Slug).
					WithCode(errx.NotFound)
			}
			return errx.Wrap(err, "failed to get user by slug", "slug", input.Slug).
				WithCode(errx.Internal)
		}

		methods, err = uc.paymentMethodRepo.WithTx(tx).
			ListByUserID(ctx, endUser.UserID, repository.UserPaymentMethodPreloadQRCodeS3Object())
		if err != nil {
			return errx.Wrap(err, "failed to list payment methods", "user_id", endUser.UserID).
				WithCode(errx.Internal)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return GetUserProfileOutput{}, err
	}

	return GetUserProfileOutput{
		User:           endUser,
		PaymentMethods: methods,
	}, nil
}
