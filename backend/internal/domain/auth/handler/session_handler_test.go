package handler_test

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/mickamy/contest"
	"github.com/stretchr/testify/require"

	authv1 "github.com/mickamy/sampay/gen/auth/v1"
	"github.com/mickamy/sampay/gen/auth/v1/authv1connect"
	"github.com/mickamy/sampay/internal/api/interceptor"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/auth/handler"
	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/domain/auth/repository"
	"github.com/mickamy/sampay/internal/infra/storage/kvs"
	"github.com/mickamy/sampay/internal/lib/cookie"
	"github.com/mickamy/sampay/internal/lib/jwt"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

func TestSession_RefreshToken(t *testing.T) {
	t.Parallel()

	userID := ulid.New()
	session := model.MustNewSession(userID)

	tests := []struct {
		name    string
		arrange func(*testing.T, *kvs.KVS) jwt.Token
		assert  func(*testing.T, *contest.Client)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, kvs *kvs.KVS) jwt.Token {
				require.NoError(t, repository.NewSession(kvs).Create(t.Context(), session))
				return session.Tokens.Refresh
			},
			assert: func(t *testing.T, ct *contest.Client) {
				ct.ExpectStatus(http.StatusOK)
			},
		},
		{
			name: "fail (refresh token not set)",
			arrange: func(t *testing.T, kvs *kvs.KVS) jwt.Token {
				require.NoError(t, repository.NewSession(kvs).Create(t.Context(), session))
				return jwt.Token{Value: "", ExpiresAt: time.Now().Add(time.Hour)}
			},
			assert: func(t *testing.T, ct *contest.Client) {
				ct.ExpectStatus(http.StatusBadRequest)
			},
		},
		{
			name: "fail (invalid refresh token)",
			arrange: func(t *testing.T, kvs *kvs.KVS) jwt.Token {
				require.NoError(t, repository.NewSession(kvs).Create(t.Context(), session))
				return jwt.Token{Value: session.Tokens.Refresh.Value + "_invalid", ExpiresAt: session.Tokens.Refresh.ExpiresAt}
			},
			assert: func(t *testing.T, ct *contest.Client) {
				ct.ExpectStatus(http.StatusBadRequest)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			infra := newInfra(t)
			refreshToken := tt.arrange(t, infra.KVS)
			c := cookie.Build("refresh_token", refreshToken.Value, refreshToken.ExpiresAt)

			// act
			ct := contest.NewWith(t,
				contest.Bind(authv1connect.NewSessionServiceHandler)(handler.NewSession(infra)),
				connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
			).
				Procedure(authv1connect.SessionServiceRefreshTokenProcedure).
				In(&authv1.RefreshTokenRequest{}).
				Header("Cookie", c.Name+"="+c.Value).
				Do()

			// assert
			tt.assert(t, ct)
		})
	}
}

func TestSession_Logout(t *testing.T) {
	t.Parallel()

	userID := ulid.New()
	session := model.MustNewSession(userID)

	tests := []struct {
		name    string
		arrange func(*testing.T, *di.Infra) (jwt.Token, jwt.Token)
		assert  func(*testing.T, *contest.Client)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, infra *di.Infra) (jwt.Token, jwt.Token) {
				require.NoError(t, repository.NewSession(infra.KVS).Create(t.Context(), session))
				return session.Tokens.Access, session.Tokens.Refresh
			},
			assert: func(t *testing.T, ct *contest.Client) {
				ct.ExpectStatus(http.StatusOK)
			},
		},
		{
			name: "fail (access token not set - rejected by auth interceptor)",
			arrange: func(t *testing.T, infra *di.Infra) (jwt.Token, jwt.Token) {
				require.NoError(t, repository.NewSession(infra.KVS).Create(t.Context(), session))
				return jwt.Token{Value: "", ExpiresAt: time.Now().Add(time.Hour)}, session.Tokens.Refresh
			},
			assert: func(t *testing.T, ct *contest.Client) {
				ct.ExpectStatus(http.StatusUnauthorized)
			},
		},
		{
			name: "fail (refresh token not set)",
			arrange: func(t *testing.T, infra *di.Infra) (jwt.Token, jwt.Token) {
				require.NoError(t, repository.NewSession(infra.KVS).Create(t.Context(), session))
				return session.Tokens.Access, jwt.Token{Value: "", ExpiresAt: time.Now().Add(time.Hour)}
			},
			assert: func(t *testing.T, ct *contest.Client) {
				t.Logf("refresh token not set should return bad request: %+v", ct.Err())
				ct.ExpectStatus(http.StatusBadRequest)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			infra := newInfra(t)
			accessToken, refreshToken := tt.arrange(t, infra)

			// act
			ct := contest.NewWith(t,
				contest.Bind(authv1connect.NewSessionServiceHandler)(handler.NewSession(infra)),
				connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
			).
				Procedure(authv1connect.SessionServiceLogoutProcedure).
				In(&authv1.LogoutRequest{}).
				Header("Authorization", "Bearer "+accessToken.Value).
				Header("Cookie",
					strings.Join(
						[]string{
							cookie.Build("access_token", accessToken.Value, accessToken.ExpiresAt).String(),
							cookie.Build("refresh_token", refreshToken.Value, refreshToken.ExpiresAt).String(),
						}, "; ",
					),
				).
				Do()

			// assert
			tt.assert(t, ct)
		})
	}
}
