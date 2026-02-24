package usecase

import (
	"context"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/lib/slicex"
	"github.com/mickamy/sampay/internal/lib/ulid"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

type SavePaymentMethodsInput struct {
	PaymentMethods []SavePaymentMethodInput
}

type SavePaymentMethodInput struct {
	Type             string  `map:"Type"`
	URL              string  `map:"Url"`
	QRCodeS3ObjectID *string `map:"QrCodeS3ObjectId"`
	DisplayOrder     int     `map:"DisplayOrder"`
}

type SavePaymentMethodsOutput struct {
	PaymentMethods []model.UserPaymentMethod
}

type SavePaymentMethods interface {
	Do(ctx context.Context, input SavePaymentMethodsInput) (SavePaymentMethodsOutput, error)
}

type savePaymentMethods struct {
	_                 SavePaymentMethods           `inject:"returns"`
	_                 *di.Infra                    `inject:"param"`
	writer            *database.Writer             `inject:""`
	paymentMethodRepo repository.UserPaymentMethod `inject:""`
}

func (uc *savePaymentMethods) Do(ctx context.Context, input SavePaymentMethodsInput) (SavePaymentMethodsOutput, error) {
	userID := contexts.MustAuthenticatedUserID(ctx)

	if err := savePaymentMethodInputs(input.PaymentMethods).validate(); err != nil {
		return SavePaymentMethodsOutput{}, err
	}

	methods := slicex.Map(input.PaymentMethods, func(i SavePaymentMethodInput) *model.UserPaymentMethod {
		return &model.UserPaymentMethod{
			ID:               ulid.New(),
			UserID:           userID,
			Type:             i.Type,
			URL:              i.URL,
			QRCodeS3ObjectID: i.QRCodeS3ObjectID,
			DisplayOrder:     i.DisplayOrder,
		}
	})

	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		if err := uc.paymentMethodRepo.WithTx(tx).DeleteByUserID(ctx, userID); err != nil {
			return errx.Wrap(err, "failed to delete existing payment methods")
		}
		if len(methods) > 0 {
			if err := uc.paymentMethodRepo.WithTx(tx).CreateAll(ctx, methods); err != nil {
				return errx.Wrap(err, "failed to create payment methods")
			}
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return SavePaymentMethodsOutput{}, err
	}

	saved, err := uc.paymentMethodRepo.ListByUserID(ctx, userID)
	if err != nil {
		return SavePaymentMethodsOutput{}, errx.Wrap(err, "failed to list saved payment methods")
	}

	return SavePaymentMethodsOutput{PaymentMethods: saved}, nil
}

type savePaymentMethodInputs []SavePaymentMethodInput

func (i savePaymentMethodInputs) validate() error {
	seen := make(map[string]bool, len(i))
	for _, item := range i {
		if item.Type == "" {
			return errx.New("payment method type is required").WithCode(errx.InvalidArgument)
		}
		if item.URL == "" {
			return errx.New("payment method URL is required").WithCode(errx.InvalidArgument)
		}
		if seen[item.Type] {
			return errx.New("duplicate payment method type", "type", item.Type).WithCode(errx.InvalidArgument)
		}
		seen[item.Type] = true
	}
	return nil
}
