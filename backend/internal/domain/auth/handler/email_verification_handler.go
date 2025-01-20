package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	authModel "mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/misc/i18n"
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
	lang := contexts.MustLanguage(ctx)
	var intentType authModel.EmailVerificationIntentType
	switch req.Msg.IntentType {
	case authv1.RequestVerificationRequest_INTENT_TYPE_SIGN_UP:
		intentType = authModel.EmailVerificationIntentTypeSignUp
	case authv1.RequestVerificationRequest_INTENT_TYPE_RESET_PASSWORD:
		intentType = authModel.EmailVerificationIntentTypeResetPassword
	default:
		slogger.ErrorCtx(ctx, "invalid intent type", "intent_type", req.Msg.IntentType)
		return nil, commonResponse.
			NewError(connect.CodeInvalidArgument, errors.New("invalid intent type")).
			WithFieldViolation("intent_type", i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.AuthHandlerEmail_verificationRequestErrorInvalid_intent_type})).
			AsConnectError()
	}
	out, err := h.request.Do(ctx, usecase.RequestEmailVerificationInput{
		IntentType: intentType,
		Email:      req.Msg.Email,
	})
	if err != nil {
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			if errors.Is(err, usecase.ErrRequestEmailVerificationEmailAlreadyExists) || errors.Is(err, usecase.ErrRequestEmailVerificationAuthenticationNotFound) {
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
	out, err := h.verify.Do(ctx, usecase.VerifyEmailInput{
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
		Token: out.Token,
	})
	return res, nil
}

var _ authv1connect.EmailVerificationServiceHandler = (*EmailVerification)(nil)
