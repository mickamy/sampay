package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/config"
	v1 "github.com/mickamy/sampay/gen/user/v1"
	"github.com/mickamy/sampay/gen/user/v1/userv1connect"
	"github.com/mickamy/sampay/internal/di"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/domain/user/mapper"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/usecase"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/lib/slicex"
)

var _ userv1connect.PaymentMethodServiceHandler = (*PaymentMethod)(nil)

type PaymentMethod struct {
	_                  *di.Infra                  `inject:"param"`
	listPaymentMethods usecase.ListPaymentMethods `inject:""`
	savePaymentMethods usecase.SavePaymentMethods `inject:""`
}

func (h *PaymentMethod) ListPaymentMethods(
	ctx context.Context, _ *connect.Request[v1.ListPaymentMethodsRequest],
) (*connect.Response[v1.ListPaymentMethodsResponse], error) {
	out, err := h.listPaymentMethods.Do(ctx, usecase.ListPaymentMethodsInput{})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	cloudfrontURL := config.AWS().CloudfrontURL()
	methods := slicex.Map(out.PaymentMethods, func(m model.UserPaymentMethod) *v1.PaymentMethod {
		return mapper.ToV1PaymentMethod(m, cloudfrontURL)
	})

	return connect.NewResponse(&v1.ListPaymentMethodsResponse{
		PaymentMethods: methods,
	}), nil
}

func (h *PaymentMethod) SavePaymentMethods(
	ctx context.Context, r *connect.Request[v1.SavePaymentMethodsRequest],
) (*connect.Response[v1.SavePaymentMethodsResponse], error) {
	var inputs []usecase.SavePaymentMethodInput
	for _, pm := range r.Msg.GetPaymentMethods() {
		input, err := mapper.ToUsecaseSavePaymentMethodInputPtr(pm)
		if err != nil {
			var localizable *cmodel.LocalizableError
			if errors.As(err, &localizable) {
				return nil, errx.Wrap(err).
					WithFieldViolation("type", localizable.LocalizeContext(ctx))
			}
			return nil, errx.Wrap(err).
				WithFieldViolation("type", err.Error())
		}
		inputs = append(inputs, *input)
	}

	out, err := h.savePaymentMethods.Do(ctx, usecase.SavePaymentMethodsInput{
		PaymentMethods: inputs,
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	cloudfrontURL := config.AWS().CloudfrontURL()
	methods := slicex.Map(out.PaymentMethods, func(m model.UserPaymentMethod) *v1.PaymentMethod {
		return mapper.ToV1PaymentMethod(m, cloudfrontURL)
	})

	return connect.NewResponse(&v1.SavePaymentMethodsResponse{
		PaymentMethods: methods,
	}), nil
}
