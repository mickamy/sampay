package handler

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	"github.com/mickamy/errx"

	authv1 "github.com/mickamy/sampay/gen/auth/v1"
	"github.com/mickamy/sampay/gen/auth/v1/authv1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/auth/usecase"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	cresponse "github.com/mickamy/sampay/internal/domain/common/response"
	"github.com/mickamy/sampay/internal/lib/cookie"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/misc/i18n"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var _ authv1connect.SessionServiceHandler = (*Session)(nil)

type Session struct {
	_            *di.Infra            `inject:"param"`
	refreshToken usecase.RefreshToken `inject:""`
	logout       usecase.Logout       `inject:""`
}

func (h *Session) RefreshToken(
	ctx context.Context, r *connect.Request[authv1.RefreshTokenRequest],
) (*connect.Response[authv1.RefreshTokenResponse], error) {
	cookieInHeader := r.Header().Get("cookie")
	refreshToken, err := cookie.ParseRefreshToken(cookieInHeader)
	if err != nil {
		return nil, errx.Wrap(err).
			WithCode(errx.InvalidArgument).
			WithDetails(errx.FieldViolation(
				"refresh_token",
				i18n.LocalizeContext(ctx, messages.AuthHandlerRefreshTokenTokenNotSet()),
			))
	}

	out, err := h.refreshToken.Do(ctx, usecase.RefreshTokenInput{
		Token: refreshToken,
	})
	if err != nil {
		var localizable *cmodel.LocalizableError
		if errors.As(err, &localizable) {
			return nil, errx.Wrap(err).
				WithCode(errx.InvalidArgument).
				WithDetails(errx.FieldViolation("refresh_token", localizable.LocalizeContext(ctx)))
		}

		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, cresponse.NewInternalErrorContext(ctx, err).AsConnectError()
	}

	res := connect.NewResponse(&authv1.RefreshTokenResponse{})
	atCookie := cookie.Build(
		"access_token", out.Tokens.Access.Value,
		out.Tokens.Access.ExpiresAt,
	)
	rtCookie := cookie.Build(
		"refresh_token", out.Tokens.Refresh.Value,
		out.Tokens.Refresh.ExpiresAt,
	)
	res.Header().Add("Set-Cookie", atCookie.String())
	res.Header().Add("Set-Cookie", rtCookie.String())

	return res, nil
}

func (h *Session) Logout(
	ctx context.Context, r *connect.Request[authv1.LogoutRequest],
) (*connect.Response[authv1.LogoutResponse], error) {
	c := r.Header().Get("cookie")
	at, err := cookie.ParseAccessToken(c)
	if err != nil {
		logger.Warn(ctx, "failed to parse access token", "err", err, "cookie", c)
		return nil, errx.Wrap(err).
			WithCode(errx.InvalidArgument).
			WithDetails(
				errx.FieldViolation(
					"access_token",
					i18n.LocalizeContext(ctx, messages.AuthHandlerLogoutAccessTokenNotSet()),
				),
			)
	}
	rt, err := cookie.ParseRefreshToken(c)
	if err != nil {
		logger.Warn(ctx, "failed to parse refresh token", "err", err, "cookie", c)
		return nil, errx.Wrap(err).
			WithCode(errx.InvalidArgument).
			WithDetails(
				errx.FieldViolation(
					"refresh_token",
					i18n.LocalizeContext(ctx, messages.AuthHandlerLogoutRefreshTokenNotSet()),
				),
			)
	}

	_, err = h.logout.Do(ctx, usecase.LogoutInput{
		AccessToken:  at,
		RefreshToken: rt,
	})
	if err != nil {
		var localizable *cmodel.LocalizableError
		if errors.As(err, &localizable) {
			var field string
			switch {
			case errors.Is(localizable, usecase.ErrLogoutInvalidAccessToken):
				field = "access_token"
			case errors.Is(localizable, usecase.ErrLogoutInvalidRefreshToken):
				field = "refresh_token"
			case errors.Is(localizable, usecase.ErrLogoutTokenMismatch):
				field = "access_token"
			default:
				field = "unknown"
				logger.Error(ctx, "unexpected localizable error", "err", err)
			}
			return nil, errx.Wrap(err).WithCode(errx.InvalidArgument).
				WithDetails(errx.FieldViolation(field, localizable.LocalizeContext(ctx)))
		}

		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, cresponse.NewInternalErrorContext(ctx, err).AsConnectError()
	}

	res := connect.NewResponse(&authv1.LogoutResponse{})
	expiredAt := time.Now().Add(-time.Hour)
	res.Header().Add("Set-Cookie", cookie.Build("access_token", "", expiredAt).String())
	res.Header().Add("Set-Cookie", cookie.Build("refresh_token", "", expiredAt).String())
	return res, nil
}
