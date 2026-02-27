package usecase

import (
	"context"
	"errors"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/repository"
	umodel "github.com/mickamy/sampay/internal/domain/user/model"
	urepository "github.com/mickamy/sampay/internal/domain/user/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var ErrGetEventNotFound = cmodel.NewLocalizableError(
	errx.NewSentinel("event not found", errx.NotFound),
).WithMessages(messages.EventUseCaseErrorNotFound())

type GetEventInput struct {
	ID string
}

type GetEventOutput struct {
	Event          model.Event
	User           umodel.EndUser
	PaymentMethods []umodel.UserPaymentMethod
}

type GetEvent interface {
	Do(ctx context.Context, input GetEventInput) (GetEventOutput, error)
}

type getEvent struct {
	_                 GetEvent                      `inject:"returns"`
	_                 *di.Infra                     `inject:"param"`
	reader            *database.Reader              `inject:""`
	eventRepo         repository.Event              `inject:""`
	endUserRepo       urepository.EndUser           `inject:""`
	paymentMethodRepo urepository.UserPaymentMethod `inject:""`
}

func (uc *getEvent) Do(ctx context.Context, input GetEventInput) (GetEventOutput, error) {
	var (
		ev      model.Event
		endUser umodel.EndUser
		methods []umodel.UserPaymentMethod
	)

	if err := uc.reader.Transaction(ctx, func(tx *database.DB) error {
		var err error
		ev, err = uc.eventRepo.WithTx(tx).Get(ctx, input.ID, repository.EventPreloadParticipants())
		if err != nil {
			if errors.Is(err, database.ErrNotFound) {
				return ErrGetEventNotFound
			}
			return errx.Wrap(err, "message", "failed to get event", "id", input.ID).
				WithCode(errx.Internal)
		}

		endUser, err = uc.endUserRepo.WithTx(tx).Get(ctx, ev.UserID)
		if err != nil {
			return errx.Wrap(err, "message", "failed to get user", "user_id", ev.UserID).
				WithCode(errx.Internal)
		}

		methods, err = uc.paymentMethodRepo.WithTx(tx).
			ListByUserID(ctx, ev.UserID, urepository.UserPaymentMethodPreloadQRCodeS3Object())
		if err != nil {
			return errx.Wrap(err, "message", "failed to list payment methods", "user_id", ev.UserID).
				WithCode(errx.Internal)
		}

		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return GetEventOutput{}, err
	}

	model.CalcAmounts(ev.TotalAmount, ev.Participants)

	return GetEventOutput{
		Event:          ev,
		User:           endUser,
		PaymentMethods: methods,
	}, nil
}
