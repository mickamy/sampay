package interceptor

import (
	"context"
	"net/http"
	"time"

	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	"github.com/bufbuild/connect-go"
)

func Cookie() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			res, err := next(ctx, req)
			if err != nil {
				return nil, err
			}

			var tokens *authv1.Tokens
			switch typedRes := res.Any().(type) {
			case *authv1.SignInResponse:
				tokens = typedRes.Tokens
			case *authv1.RefreshResponse:
				tokens = typedRes.Tokens
			case *authv1.SignOutResponse:
				removeSessionCookies(res)
				return res, nil
			default:
				return res, nil
			}
			setSessionCookies(tokens, res)
			return res, nil
		}
	}
}

func setSessionCookies(tokens *authv1.Tokens, res connect.AnyResponse) {
	at := buildCookie("access_token", tokens.Access.Value, tokens.Access.ExpiresAt.AsTime())
	res.Header().Add("Set-Cookie", at.String())
	rt := buildCookie("refresh_token", tokens.Refresh.Value, tokens.Refresh.ExpiresAt.AsTime())
	res.Header().Add("Set-Cookie", rt.String())
}

func removeSessionCookies(res connect.AnyResponse) {
	at := buildCookie("access_token", "", time.Unix(0, 0))
	res.Header().Add("Set-Cookie", at.String())
	rt := buildCookie("refresh_token", "", time.Unix(0, 0))
	res.Header().Add("Set-Cookie", rt.String())
}

func buildCookie(name, value string, expires time.Time) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Expires = expires
	cookie.HttpOnly = true
	cookie.Secure = true
	return cookie
}
