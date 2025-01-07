package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/registration/v1/registrationv1connect"
	registrationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/registration/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	authDTO "mickamy.com/sampay/internal/domain/auth/dto"
	dto "mickamy.com/sampay/internal/domain/common/dto"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/misc/i18n"
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
		if errors.Is(err, usecase.ErrCreateAccountEmailAlreadyExists) {
			return nil, dto.NewBadRequest(err).
				WithFieldViolation("email", i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: "registration.handler.error.email_already_exists"})).
				AsConnectError()
		}
		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.SignUpResponse{
		UserId: out.Session.UserID,
		Tokens: authDTO.NewTokens(out.Session.Tokens),
	})
	return res, nil
}

var _ registrationv1connect.AccountServiceHandler = (*Account)(nil)
