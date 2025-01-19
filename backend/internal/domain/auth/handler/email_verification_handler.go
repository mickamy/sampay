package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	"github.com/bufbuild/connect-go"
	"github.com/mickamy/slogger"

	authResponse "mickamy.com/sampay/internal/domain/auth/dto/response"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/lib/contexts"
)

type EmailVerification struct {
	request usecase.RequestEmailVerification
	verify  usecase.VerifyEmail
}

func NewEmailVerification(
	request usecase.RequestEmailVerification,
	verify usecase.VerifyEmail,
) *EmailVerification {
	return &EmailVerification{
		request: request,
		verify:  verify,
	}
}

func (h *EmailVerification) RequestVerification(
	ctx context.Context,
	req *connect.Request[authv1.RequestVerificationRequest],
) (*connect.Response[authv1.RequestVerificationResponse], error) {
	out, err := h.request.Do(ctx, usecase.RequestEmailVerificationInput{
		Email: req.Msg.Email,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			if errors.Is(err, usecase.ErrRequestEmailVerificationEmailAlreadyExists) {
				return nil, localizable.AsFieldViolations("email").AsConnectError()
			}
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&authv1.RequestVerificationResponse{
		Token: out.Token,
	})
	return res, nil
}

func (h *EmailVerification) VerifyEmail(
	ctx context.Context,
	req *connect.Request[authv1.VerifyEmailRequest],
) (*connect.Response[authv1.VerifyEmailResponse], error) {
	got, err := h.verify.Do(ctx, usecase.VerifyEmailInput{
		Token:   req.Msg.Token,
		PINCode: req.Msg.PinCode,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			if errors.Is(err, usecase.ErrVerifyEmailInvalidToken) {
				return nil, localizable.AsFieldViolations("pin_code").AsConnectError()
			}
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&authv1.VerifyEmailResponse{
		Tokens: authResponse.NewTokens(got.Session.Tokens),
	})
	return res, nil
}

var _ authv1connect.EmailVerificationServiceHandler = (*EmailVerification)(nil)
