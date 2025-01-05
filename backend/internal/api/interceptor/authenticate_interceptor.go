package interceptor

import (
	"context"
	"errors"
	"strings"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/auth/v1/authv1connect"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	"mickamy.com/sampay/internal/domain/auth/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
)

var (
	ErrNoAccessToken = errors.New("no access token found")
)

var authSkippingProcedures = []string{
	authv1connect.SessionServiceSignInProcedure,
	authv1connect.SessionServiceRefreshProcedure,
}

func skipper(req connect.AnyRequest) bool {
	for _, procedure := range authSkippingProcedures {
		if req.Spec().Procedure == procedure {
			return true
		}
	}
	return false
}

func Authenticate(uc usecase.AuthenticateUser) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if skipper(req) {
				return next(ctx, req)
			}

			accessToken := extractAccessToken(ctx, req)
			if accessToken == "" {
				slogger.WarnCtx(ctx, "no access token found")
				return nil, connect.NewError(connect.CodeUnauthenticated, ErrNoAccessToken)
			}
			out, err := uc.Do(ctx, usecase.AuthenticateUserInput{
				AccessToken: accessToken,
			})
			if err != nil {
				slogger.WarnCtx(ctx, "failed to authenticate user", "err", err)
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			ctx = contexts.SetAuthenticatedUser(ctx, out.User)
			return next(ctx, req)
		}
	}
}

func extractAccessToken(ctx context.Context, req connect.AnyRequest) string {
	authHeader := req.Header().Get("Authorization")
	if authHeader != "" {
		slogger.DebugCtx(ctx, "extracting access token from Authorization header")
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	cookie := req.Header().Get("Cookie")
	if cookie == "" {
		slogger.DebugCtx(ctx, "no access token found in Authorization header or Cookie")
		return ""
	}
	parts := strings.Split(cookie, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "access_token=") {
			slogger.DebugCtx(ctx, "extracting access token from Cookie")
			return strings.TrimPrefix(part, "access_token=")
		}
	}
	slogger.DebugCtx(ctx, "no access token found in Cookie")
	return ""
}
