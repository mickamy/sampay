package handler

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/mickamy/errx"

	v1 "github.com/mickamy/sampay/gen/auth/v1"
	"github.com/mickamy/sampay/gen/auth/v1/authv1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/auth/mapper"
	"github.com/mickamy/sampay/internal/domain/auth/usecase"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	cresponse "github.com/mickamy/sampay/internal/domain/common/response"
	umapper "github.com/mickamy/sampay/internal/domain/user/mapper"
	"github.com/mickamy/sampay/internal/lib/cookie"
	"github.com/mickamy/sampay/internal/lib/logger"
)

var _ authv1connect.OAuthServiceHandler = (*OAuth)(nil)

type OAuth struct {
	_             *di.Infra             `inject:"param"`
	getOAuthURL   usecase.GetOAuthURL   `inject:""`
	oauthCallback usecase.OAuthCallback `inject:""`
}

func (h *OAuth) GetOAuthURL(
	ctx context.Context, r *connect.Request[v1.GetOAuthURLRequest],
) (*connect.Response[v1.GetOAuthURLResponse], error) {
	provider, err := mapper.ToOAuthProvider(r.Msg.GetProvider())
	if err != nil {
		var localizable *cmodel.LocalizableError
		if errors.As(err, &localizable) {
			return nil, errx.Wrap(err).
				WithDetails(errx.FieldViolation("provider", localizable.LocalizeContext(ctx)))
		}
		return nil, errx.Wrap(err).
			WithDetails(errx.FieldViolation("provider", err.Error()))
	}

	out, err := h.getOAuthURL.Do(ctx, usecase.GetOAuthURLInput{Provider: provider})
	if err != nil {
		var localizable *cmodel.LocalizableError
		if errors.As(err, &localizable) {
			return nil, errx.Wrap(err).
				WithDetails(errx.FieldViolation("provider", localizable.LocalizeContext(ctx)))
		}

		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, cresponse.NewInternalErrorContext(ctx, err).AsConnectError()
	}

	return connect.NewResponse(&v1.GetOAuthURLResponse{Url: out.AuthenticationURL}), nil
}

func (h *OAuth) OAuthCallback(
	ctx context.Context, r *connect.Request[v1.OAuthCallbackRequest],
) (*connect.Response[v1.OAuthCallbackResponse], error) {
	provider, err := mapper.ToOAuthProvider(r.Msg.GetProvider())
	if err != nil {
		var localizable *cmodel.LocalizableError
		if errors.As(err, &localizable) {
			return nil, errx.Wrap(err).
				WithDetails(errx.FieldViolation("provider", localizable.LocalizeContext(ctx)))
		}
		return nil, errx.Wrap(err).
			WithDetails(errx.FieldViolation("provider", err.Error()))
	}

	out, err := h.oauthCallback.Do(ctx, usecase.OAuthCallbackInput{
		Provider: provider,
		Code:     r.Msg.GetCode(),
	})
	if err != nil {
		var localizable *cmodel.LocalizableError
		if errors.As(err, &localizable) {
			if errors.Is(localizable, usecase.ErrOAuthCallbackUnsupportedProvider) {
				return nil, errx.Wrap(err).
					WithDetails(errx.FieldViolation("provider", localizable.LocalizeContext(ctx)))
			}
			return nil, errx.Wrap(err).
				WithDetails(errx.FieldViolation("code", localizable.LocalizeContext(ctx)))
		}

		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, cresponse.NewInternalErrorContext(ctx, err).AsConnectError()
	}

	res := connect.NewResponse(&v1.OAuthCallbackResponse{
		User: umapper.ToV1UserPtr(&out.EndUser),
	})
	atCookie := cookie.Build(
		"access_token", out.Session.Tokens.Access.Value,
		out.Session.Tokens.Access.ExpiresAt,
	)
	rtCookie := cookie.Build(
		"refresh_token", out.Session.Tokens.Refresh.Value,
		out.Session.Tokens.Refresh.ExpiresAt,
	)
	res.Header().Add("Set-Cookie", atCookie.String())
	res.Header().Add("Set-Cookie", rtCookie.String())

	return res, nil
}
