package interceptor

import (
	"context"
	"errors"
	"strings"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/auth/v1/authv1connect"
	"buf.build/gen/go/mickamy/sampay/connectrpc/go/message/v1/messagev1connect"
	"buf.build/gen/go/mickamy/sampay/connectrpc/go/oauth/v1/oauthv1connect"
	"buf.build/gen/go/mickamy/sampay/connectrpc/go/registration/v1/registrationv1connect"
	"buf.build/gen/go/mickamy/sampay/connectrpc/go/user/v1/userv1connect"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	"mickamy.com/sampay/internal/domain/auth/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
)

var (
	ErrNoBearerToken = errors.New("no bearer token found")
)

var authSkippingProcedures = []string{
	authv1connect.EmailVerificationServiceRequestVerificationProcedure,
	authv1connect.EmailVerificationServiceVerifyEmailProcedure,
	authv1connect.SessionServiceSignInProcedure,
	authv1connect.SessionServiceRefreshProcedure,
	messagev1connect.MessageServiceSendMessageProcedure,
	oauthv1connect.OAuthServiceSignInProcedure,
	oauthv1connect.OAuthServiceGoogleCallbackProcedure,
	userv1connect.UserServiceGetUserProcedure,
}

var anonymousProcedures = []string{
	registrationv1connect.OnboardingServiceGetOnboardingStepProcedure,
	registrationv1connect.OnboardingServiceCreatePasswordProcedure,
}

func authSkippingPaths(req connect.AnyRequest) bool {
	for _, procedure := range authSkippingProcedures {
		if req.Spec().Procedure == procedure {
			return true
		}
	}
	return false
}

func anonymousPaths(req connect.AnyRequest) bool {
	for _, procedure := range anonymousProcedures {
		if req.Spec().Procedure == procedure {
			return true
		}
	}
	return false
}

func Authenticate(authenticate usecase.AuthenticateUser, anonymous usecase.AuthenticateAnonymousUser) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if authSkippingPaths(req) {
				return next(ctx, req)
			}

			bearer := extractBearerToken(ctx, req)
			if bearer == "" {
				slogger.WarnCtx(ctx, "no bearer token found")
				return nil, connect.NewError(connect.CodeUnauthenticated, ErrNoBearerToken)
			}

			out, err := authenticate.Do(ctx, usecase.AuthenticateUserInput{
				Token: bearer,
			})
			if err != nil {
				slogger.WarnCtx(ctx, "failed to authenticate user", "err", err)
				return nil, connect.NewError(connect.CodeUnauthenticated, err)
			}
			if out.User == nil {
				if anonymousPaths(req) {
					_, err := anonymous.Do(ctx, usecase.AuthenticateAnonymousUserInput{Token: bearer})
					if err != nil {
						slogger.WarnCtx(ctx, "failed to authenticate anonymous user", "err", err)
						return nil, connect.NewError(connect.CodeUnauthenticated, err)
					}
					ctx = contexts.SetAnonymousUserToken(ctx, bearer)
					return next(ctx, req)
				} else {
					return nil, connect.NewError(connect.CodeUnauthenticated, usecase.ErrAuthenticateUserUserNotFound)
				}
			} else {
				ctx = contexts.SetAuthenticatedUserID(ctx, out.User.ID)
			}
			return next(ctx, req)
		}
	}
}

func extractBearerToken(ctx context.Context, req connect.AnyRequest) string {
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
