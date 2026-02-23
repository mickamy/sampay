package interceptor

import (
	"context"
	"errors"
	"slices"
	"strings"

	"connectrpc.com/connect"
	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/gen/auth/v1/authv1connect"
	ausecase "github.com/mickamy/sampay/internal/domain/auth/usecase"

	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/lib/cookie"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

var authSkippingProcedures = []string{
	authv1connect.OAuthServiceGetOAuthURLProcedure,
	authv1connect.OAuthServiceOAuthCallbackProcedure,
	authv1connect.SessionServiceRefreshTokenProcedure,
}

func Authenticate(uc ausecase.Authenticate) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if authSkippingPaths(req) {
				return next(ctx, req)
			}

			token := extractBearerToken(req)
			out, err := uc.Do(ctx, ausecase.AuthenticateInput{
				Token: token,
			})

			if err != nil {
				var localizable *cmodel.LocalizableError
				if errors.As(err, &localizable) {
					return nil, errx.Wrap(localizable).
						WithFieldViolation("access_token", localizable.LocalizeContext(ctx))
				}

				logger.Error(ctx, "failed to execute use-case", "err", err)
				return nil, connect.NewError(connect.CodeInternal, err)
			}

			ctx = contexts.SetAuthenticatedUserID(ctx, out.UserID)

			return next(ctx, req)
		}
	}
}

func authSkippingPaths(req connect.AnyRequest) bool {
	return slices.Contains(authSkippingProcedures, req.Spec().Procedure)
}

func extractBearerToken(req connect.AnyRequest) string {
	// check for Authorization header first
	authHeader := strings.TrimSpace(req.Header().Get("Authorization"))
	if authHeader != "" {
		parts := strings.Fields(authHeader)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") && parts[1] != "" {
			return parts[1]
		}
	}

	// if Authorization header is not present or not a valid Bearer token, check for access token in Cookie
	token, err := cookie.ParseAccessToken(req.Header().Get("Cookie"))
	if err != nil {
		return ""
	}

	return token
}
