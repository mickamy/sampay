package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	authResponse "mickamy.com/sampay/internal/domain/auth/dto/response"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
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
	delete usecase.DeleteSession,
) *Session {
	return &Session{
		create:  create,
		refresh: refresh,
		delete:  delete,
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
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&authv1.SignInResponse{
		UserId: out.Session.UserID,
		Tokens: authResponse.NewTokens(out.Session.Tokens),
	})
	return res, nil
}

func (h *Session) Refresh(
	ctx context.Context,
	req *connect.Request[authv1.RefreshRequest],
) (*connect.Response[authv1.RefreshResponse], error) {
	lang := contexts.MustLanguage(ctx)

	tkn := req.Msg.RefreshToken
	if tkn == nil {
		tknFromCookie, err := extractRefreshTokenFromCookie(req)
		if err != nil {
			slogger.ErrorCtx(ctx, "failed to extract refresh token from cookie", "err", err)
			return nil, commonResponse.NewBadRequest(err).
				WithMessage(i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.AuthUsecaseErrorInvalid_refresh_token})).
				AsConnectError()
		}
		tkn = &tknFromCookie
	}

	out, err := h.refresh.Do(ctx, usecase.RefreshSessionInput{
		RefreshToken: *tkn,
	})
	if err != nil {
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}

	return connect.NewResponse(&authv1.RefreshResponse{
		Tokens: authResponse.NewTokens(out.Tokens),
	}), nil
}

func extractRefreshTokenFromCookie(req *connect.Request[authv1.RefreshRequest]) (string, error) {
	cookies, err := http.ParseCookie(req.Header().Get("Cookie"))
	if err != nil {
		return "", fmt.Errorf("failed to parse cookie: %w", err)
	}
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("refresh token not found")
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
		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		// do not return error if deleting tokens failed
		if !errors.Is(err, usecase.ErrDeleteSessionDeletingTokensFailed) {
			slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
			return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
		}
	}

	return connect.NewResponse(&authv1.SignOutResponse{}), nil
}

var _ authv1connect.SessionServiceHandler = (*Session)(nil)
