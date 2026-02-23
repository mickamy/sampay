package handler_test

import (
	"net/http"
	"testing"

	"connectrpc.com/connect"
	"github.com/mickamy/contest"
	"github.com/stretchr/testify/assert"

	authv1 "github.com/mickamy/sampay/gen/auth/v1"
	"github.com/mickamy/sampay/gen/auth/v1/authv1connect"
	"github.com/mickamy/sampay/internal/api/interceptor"
	"github.com/mickamy/sampay/internal/domain/auth/handler"
)

func TestOAuth_OAuthCallback(t *testing.T) {
	t.Parallel()
	t.Skip("requires real OAuth provider interaction")
}

func TestOAuth_GetOAuthURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		provider authv1.OAuthProvider
		assert   func(*testing.T, *contest.Client)
	}{
		{
			name:     "google",
			provider: authv1.OAuthProvider_O_AUTH_PROVIDER_GOOGLE,
			assert: func(t *testing.T, ct *contest.Client) {
				var out authv1.GetOAuthURLResponse
				ct.ExpectStatus(http.StatusOK).
					Out(&out)
				assert.Contains(t, out.GetUrl(), "accounts.google.com")
			},
		},
		{
			name:     "line",
			provider: authv1.OAuthProvider_O_AUTH_PROVIDER_LINE,
			assert: func(t *testing.T, ct *contest.Client) {
				var out authv1.GetOAuthURLResponse
				ct.ExpectStatus(http.StatusOK).
					Out(&out)
				assert.Contains(t, out.GetUrl(), "access.line.me")
			},
		},
		{
			name:     "unspecified provider",
			provider: authv1.OAuthProvider_O_AUTH_PROVIDER_UNSPECIFIED,
			assert: func(t *testing.T, ct *contest.Client) {
				t.Logf("unspecified provider should return bad request: %+v", ct.Err())
				ct.ExpectStatus(http.StatusBadRequest)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			infra := newInfra(t)

			// act
			ct := contest.NewWith(t,
				contest.Bind(authv1connect.NewOAuthServiceHandler)(handler.NewOAuth(infra)),
				connect.WithInterceptors(interceptor.NewInterceptors(infra)...),
			).
				Procedure(authv1connect.OAuthServiceGetOAuthURLProcedure).
				In(&authv1.GetOAuthURLRequest{Provider: tt.provider}).
				Do()

			// assert
			tt.assert(t, ct)
		})
	}
}
