package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/domain/auth/usecase"
	"github.com/mickamy/sampay/internal/lib/oauth"
)

func TestGetOAuthURL_Do(t *testing.T) {
	t.Parallel()

	resolver := oauth.NewResolverFromConfig(config.OAuth())

	tests := []struct {
		name     string
		provider model.OAuthProvider
		assert   func(t *testing.T, got usecase.GetOAuthURLOutput, err error)
	}{
		{
			name:     "google",
			provider: model.OAuthProviderGoogle,
			assert: func(t *testing.T, got usecase.GetOAuthURLOutput, err error) {
				require.NoError(t, err)
				assert.Contains(t, got.AuthenticationURL, "accounts.google.com")
			},
		},
		{
			name:     "line",
			provider: model.OAuthProviderLINE,
			assert: func(t *testing.T, got usecase.GetOAuthURLOutput, err error) {
				require.NoError(t, err)
				assert.Contains(t, got.AuthenticationURL, "access.line.me")
			},
		},
		{
			name:     "unsupported provider",
			provider: model.OAuthProvider("unknown"),
			assert: func(t *testing.T, got usecase.GetOAuthURLOutput, err error) {
				require.Error(t, err)
				assert.Empty(t, got.AuthenticationURL)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sut := usecase.NewGetOAuthURL(resolver)
			got, err := sut.Do(t.Context(), usecase.GetOAuthURLInput{
				Provider: tt.provider,
			})

			tt.assert(t, got, err)
		})
	}
}
