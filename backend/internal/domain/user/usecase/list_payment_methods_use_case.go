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

type ListPaymentMethodsInput struct{}

type ListPaymentMethodsOutput struct {
	PaymentMethods []model.UserPaymentMethod
}

type ListPaymentMethods interface {
	Do(ctx context.Context, input ListPaymentMethodsInput) (ListPaymentMethodsOutput, error)
}

type listPaymentMethods struct {
	_                 ListPaymentMethods           `inject:"returns"`
	_                 *di.Infra                    `inject:"param"`
	reader            *database.Reader             `inject:""`
	paymentMethodRepo repository.UserPaymentMethod `inject:""`
}

func (uc *listPaymentMethods) Do(ctx context.Context, _ ListPaymentMethodsInput) (ListPaymentMethodsOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	var methods []model.UserPaymentMethod
	if err := uc.reader.Transaction(ctx, func(tx *database.DB) error {
		var err error
		methods, err = uc.paymentMethodRepo.ListByUserID(ctx, userID)
		if err != nil {
			return errx.Wrap(err, "failed to list payment methods", "user_id", userID)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return ListPaymentMethodsOutput{}, err
	}

	return ListPaymentMethodsOutput{PaymentMethods: methods}, nil
}
