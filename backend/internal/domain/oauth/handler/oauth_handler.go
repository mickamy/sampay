package handler

import (
	"context"
	"net/http"
	"time"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/oauth/v1/oauthv1connect"
	oauthv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/oauth/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/domain/oauth/dto/request"
	"mickamy.com/sampay/internal/domain/oauth/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/misc/i18n"
)

type OAuth struct {
	signIn   usecase.OAuthSignIn
	callback usecase.OAuthCallback
}

func NewOAuth(
	signIn usecase.OAuthSignIn,
	callback usecase.OAuthCallback,
) *OAuth {
	return &OAuth{
		signIn:   signIn,
		callback: callback,
	}
}

func (h *OAuth) SignIn(
	ctx context.Context,
	req *connect.Request[oauthv1.SignInRequest],
) (*connect.Response[oauthv1.SignInResponse], error) {
	lang := contexts.MustLanguage(ctx)
	provider, err := request.NewOAuthProvider(req.Msg.Provider)
	if err != nil {
		slogger.ErrorCtx(ctx, "failed to parse request", "err", err)
		return nil, commonResponse.NewBadRequest(err).
			WithFieldViolation("provider", i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.OauthHandlerErrorInvalid_provider_type})).
			AsConnectError()
	}

	out, err := h.signIn.Do(ctx, usecase.OAuthSignInInput{
		Provider: *provider,
	})
	if err != nil {
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&oauthv1.SignInResponse{
		AuthorizationUrl: out.AuthenticationURL,
	})
	return res, nil
}

func (h *OAuth) GoogleCallback(
	ctx context.Context,
	req *connect.Request[oauthv1.GoogleCallbackRequest],
) (*connect.Response[oauthv1.GoogleCallbackResponse], error) {
	out, err := h.callback.Do(ctx, usecase.OAuthCallbackInput{
		Code: req.Msg.Code,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&oauthv1.GoogleCallbackResponse{})
	res.Header().Add("Set-Cookie", newVerifiedCookie(out.VerifiedToken).String())
	return res, nil
}

func newVerifiedCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     "verified",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Second),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}
}

var _ oauthv1connect.OAuthServiceHandler = (*OAuth)(nil)
