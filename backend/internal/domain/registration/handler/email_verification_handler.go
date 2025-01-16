package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/registration/v1/registrationv1connect"
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
) *EmailVerification {
	return &EmailVerification{
		request: request,
	}
}

func (h *EmailVerification) RequestVerification(
	ctx context.Context,
	req *connect.Request[registrationv1.RequestVerificationRequest],
) (*connect.Response[registrationv1.RequestVerificationResponse], error) {
	out, err := h.request.Do(ctx, usecase.RequestEmailVerificationInput{
		Email: req.Msg.Email,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			if errors.Is(err, usecase.ErrRequestEmailVerificationEmailAlreadyExists) {
				return nil, localizable.AsFieldViolations("token").AsConnectError()
			}
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.RequestVerificationResponse{
		Token: out.Token,
	})
	return res, nil
}

func (h *EmailVerification) VerifyEmail(
	ctx context.Context,
	req *connect.Request[registrationv1.VerifyEmailRequest],
) (*connect.Response[registrationv1.VerifyEmailResponse], error) {
	_, err := h.verify.Do(ctx, usecase.VerifyEmailInput{
		Token: req.Msg.Token,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.VerifyEmailResponse{})
	return res, nil
}

var _ registrationv1connect.EmailVerificationServiceHandler = (*EmailVerification)(nil)
