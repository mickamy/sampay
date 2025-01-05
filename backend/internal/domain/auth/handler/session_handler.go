package handler

import (
	"context"
	"errors"
	"fmt"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"
	"google.golang.org/protobuf/types/known/timestamppb"

	"mickamy.com/sampay/internal/domain/auth/usecase"
	dto "mickamy.com/sampay/internal/domain/common/dto"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/jwt"
	"mickamy.com/sampay/internal/misc/i18n"
)

type Session struct {
	create  usecase.CreateSession
	refresh usecase.RefreshSession
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
		Tokens: h.newTokens(out.Session.Tokens),
	})
	return res, nil
}

func (h *Session) Refresh(
	ctx context.Context,
	req *connect.Request[authv1.RefreshRequest],
) (*connect.Response[authv1.RefreshResponse], error) {
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
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to execute use case: %w", err))
	}

	return connect.NewResponse(&authv1.RefreshResponse{
		Tokens: h.newTokens(out.Tokens),
	}), nil
}

func (h *Session) SignOut(
	ctx context.Context,
	req *connect.Request[authv1.SignOutRequest],
) (*connect.Response[authv1.SignOutResponse], error) {
	panic("implement me")
}

func (h *Session) newTokens(tokens jwt.Tokens) *authv1.Tokens {
	return &authv1.Tokens{
		Access:  h.newToken(tokens.Access),
		Refresh: h.newToken(tokens.Refresh),
	}
}

func (h *Session) newToken(token jwt.Token) *authv1.Token {
	return &authv1.Token{
		Value:     token.Value,
		ExpiresAt: timestamppb.New(token.ExpiresAt),
	}
}

var _ authv1connect.SessionServiceHandler = (*Session)(nil)
