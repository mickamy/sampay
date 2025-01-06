package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	authDTO "mickamy.com/sampay/internal/domain/auth/dto"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	dto "mickamy.com/sampay/internal/domain/common/dto"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/misc/i18n"
)

type Session struct {
	create  usecase.CreateSession
	refresh usecase.RefreshSession
	delete  usecase.DeleteSession
}

func NewSession(
	create usecase.CreateSession,
	refresh usecase.RefreshSession,
) *Session {
	return &Session{
		create:  create,
		refresh: refresh,
	}
}

func (h *Session) SignIn(
	ctx context.Context,
	req *connect.Request[authv1.SignInRequest],
) (*connect.Response[authv1.SignInResponse], error) {
	out, err := h.create.Do(ctx, usecase.CreateSessionInput{
		Email:    req.Msg.Email,
		Password: req.Msg.Password,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if errors.Is(err, usecase.ErrCreateSessionPasswordNotMatch) {
			return nil, dto.NewBadRequest(err).
				WithMessage(i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: "auth.error.invalid_email_password"})).
				AsConnectError()
		}
		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&authv1.SignInResponse{
		UserId: out.Session.UserID,
		Tokens: authDTO.NewTokens(out.Session.Tokens),
	})
	return res, nil
}

func (h *Session) Refresh(
	ctx context.Context,
	req *connect.Request[authv1.RefreshRequest],
) (*connect.Response[authv1.RefreshResponse], error) {
	// TODO: get refresh token from cookie
	out, err := h.refresh.Do(ctx, usecase.RefreshSessionInput{
		RefreshToken: *req.Msg.RefreshToken,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrRefreshSessionTokenNotFound) {
			return nil, dto.NewBadRequest(err).
				WithMessage(i18n.MustLocalizeMessageCtx(ctx, i18n.Config{MessageID: "auth.error.invalid_refresh_token"})).
				AsConnectError()
		}
		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}

	return connect.NewResponse(&authv1.RefreshResponse{
		Tokens: authDTO.NewTokens(out.Tokens),
	}), nil
}

func (h *Session) SignOut(
	ctx context.Context,
	req *connect.Request[authv1.SignOutRequest],
) (*connect.Response[authv1.SignOutResponse], error) {
	_, err := h.delete.Do(ctx, usecase.DeleteSessionInput{
		AccessToken:  req.Msg.AccessToken,
		RefreshToken: req.Msg.RefreshToken,
	})
	if err != nil {
		// do not return error if deleting tokens failed
		if !errors.Is(err, usecase.ErrDeleteSessionDeletingTokensFailed) {
			slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
			return nil, dto.NewInternalError(ctx, err).AsConnectError()
		}
	}

	return connect.NewResponse(&authv1.SignOutResponse{}), nil
}

var _ authv1connect.SessionServiceHandler = (*Session)(nil)
