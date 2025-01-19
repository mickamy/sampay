package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	"github.com/bufbuild/connect-go"
	"github.com/mickamy/slogger"

	"mickamy.com/sampay/internal/domain/auth/usecase"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/lib/contexts"
)

type PasswordReset struct {
	reset usecase.ResetPassword
}

func NewPasswordReset(
	reset usecase.ResetPassword,
) *PasswordReset {
	return &PasswordReset{
		reset: reset,
	}
}

func (h *PasswordReset) ResetPassword(
	ctx context.Context,
	req *connect.Request[authv1.ResetPasswordRequest],
) (*connect.Response[authv1.ResetPasswordResponse], error) {
	_, err := h.reset.Do(ctx, usecase.ResetPasswordInput{
		Token: req.Msg.Token,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			if errors.Is(err, usecase.ErrResetPasswordEmailVerificationInvalidToken) || errors.Is(err, usecase.ErrResetPasswordEmailVerificationAlreadyConsumed) {
				return nil, localizable.AsFieldViolations("token").AsConnectError()
			}
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&authv1.ResetPasswordResponse{})
	return res, nil
}

var _ authv1connect.PasswordResetServiceHandler = (*PasswordReset)(nil)
