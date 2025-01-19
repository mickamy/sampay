package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/registration/v1/registrationv1connect"
	registrationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/registration/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/domain/registration/usecase"
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
	req *connect.Request[registrationv1.RequestVerificationRequest],
) (*connect.Response[registrationv1.RequestVerificationResponse], error) {
	_, err := h.request.Do(ctx, usecase.RequestEmailVerificationInput{
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
	res := connect.NewResponse(&registrationv1.RequestVerificationResponse{})
	return res, nil
}

func (h *EmailVerification) VerifyEmail(
	ctx context.Context,
	req *connect.Request[registrationv1.VerifyEmailRequest],
) (*connect.Response[registrationv1.VerifyEmailResponse], error) {
	got, err := h.verify.Do(ctx, usecase.VerifyEmailInput{
		Email:   req.Msg.Email,
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
	res := connect.NewResponse(&registrationv1.VerifyEmailResponse{
		Token: got.Token,
	})
	return res, nil
}

var _ registrationv1connect.EmailVerificationServiceHandler = (*EmailVerification)(nil)
