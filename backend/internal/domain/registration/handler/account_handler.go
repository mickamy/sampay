package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/registration/v1/registrationv1connect"
	registrationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/registration/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	authResponse "mickamy.com/sampay/internal/domain/auth/dto/response"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
)

type Account struct {
	create usecase.CreateAccount
}

func NewAccount(
	create usecase.CreateAccount,
) *Account {
	return &Account{
		create: create,
	}
}

func (h *Account) SignUp(
	ctx context.Context,
	req *connect.Request[registrationv1.SignUpRequest],
) (*connect.Response[registrationv1.SignUpResponse], error) {
	out, err := h.create.Do(ctx, usecase.CreateAccountInput{
		Email:    req.Msg.Email,
		Password: req.Msg.Password,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		localizable := new(commonModel.LocalizableError)
		if errors.As(err, &localizable) {
			if errors.Is(err, usecase.ErrCreateAccountEmailAlreadyExists) {
				return nil, commonResponse.NewBadRequest(err).
					WithFieldViolation("email", localizable.Localize(lang)).
					AsConnectError()
			}
			return nil, commonResponse.NewBadRequest(err).WithMessage(localizable.Localize(lang)).AsConnectError()
		}
		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.SignUpResponse{
		UserId: out.Session.UserID,
		Tokens: authResponse.NewTokens(out.Session.Tokens),
	})
	return res, nil
}

var _ registrationv1connect.AccountServiceHandler = (*Account)(nil)
