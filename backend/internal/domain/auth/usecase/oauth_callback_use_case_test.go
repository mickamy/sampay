package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/domain/auth/query"
	"github.com/mickamy/sampay/internal/domain/auth/usecase"
	ufixture "github.com/mickamy/sampay/internal/domain/user/fixture"
	umodel "github.com/mickamy/sampay/internal/domain/user/model"
	uquery "github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/lib/oauth"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

type fakeOAuthClient struct {
	payload oauth.Payload
	err     error
}

func (c *fakeOAuthClient) AuthenticationURL() (string, error) {
	return "https://example.com/auth", nil
}

func (c *fakeOAuthClient) Callback(_ context.Context, _ string) (oauth.Payload, error) {
	return c.payload, c.err
}

func newFakeResolver(provider oauth.Provider, client oauth.Client) *oauth.Resolver {
	return &oauth.Resolver{
		Clients: map[oauth.Provider]oauth.Client{
			provider: client,
		},
	}
}

func TestOAuthCallback_Do(t *testing.T) {
	t.Parallel()

	fakeUID := ulid.New()
	fakePayload := oauth.Payload{
		Provider: oauth.ProviderGoogle,
		UID:      fakeUID,
		Name:     "Test User",
		Email:    "test@example.com",
	}

	tests := []struct {
		name    string
		arrange func(t *testing.T, infra *di.Infra) *oauth.Resolver
		input   usecase.OAuthCallbackInput
		assert  func(t *testing.T, infra *di.Infra, got usecase.OAuthCallbackOutput, err error)
	}{
		{
			name: "success (new user)",
			arrange: func(t *testing.T, infra *di.Infra) *oauth.Resolver {
				return newFakeResolver(oauth.ProviderGoogle, &fakeOAuthClient{payload: fakePayload})
			},
			input: usecase.OAuthCallbackInput{
				Provider: model.OAuthProviderGoogle,
				Code:     "valid_code",
			},
			assert: func(t *testing.T, infra *di.Infra, got usecase.OAuthCallbackOutput, err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Session.Tokens.Access.Value)
				assert.NotEmpty(t, got.Session.Tokens.Refresh.Value)
				assert.NotEmpty(t, got.EndUser)

				// verify oauth account was created
				account, err := query.OAuthAccounts(infra.DB).
					Where("provider = ? AND uid = ?", "google", fakeUID).
					First(t.Context())
				require.NoError(t, err)
				assert.Equal(t, fakeUID, account.UID)
			},
		},
		{
			name: "success (existing user)",
			arrange: func(t *testing.T, infra *di.Infra) *oauth.Resolver {
				// create existing user + oauth account
				user := ufixture.User(nil)
				endUser := ufixture.EndUser(func(m *umodel.EndUser) {
					m.UserID = user.ID
				})
				require.NoError(t, uquery.Users(infra.DB).Create(t.Context(), &user))
				require.NoError(t, uquery.EndUsers(infra.DB).Create(t.Context(), &endUser))

				existingAccount := model.OAuthAccount{
					ID:        ulid.New(),
					EndUserID: user.ID,
					Provider:  "google",
					UID:       fakeUID,
				}
				require.NoError(t, query.OAuthAccounts(infra.DB).Create(t.Context(), &existingAccount))

				return newFakeResolver(oauth.ProviderGoogle, &fakeOAuthClient{payload: fakePayload})
			},
			input: usecase.OAuthCallbackInput{
				Provider: model.OAuthProviderGoogle,
				Code:     "valid_code",
			},
			assert: func(t *testing.T, infra *di.Infra, got usecase.OAuthCallbackOutput, err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Session.Tokens.Access.Value)
				assert.NotEmpty(t, got.EndUser)
			},
		},
		{
			name: "unsupported provider",
			arrange: func(t *testing.T, infra *di.Infra) *oauth.Resolver {
				return newFakeResolver(oauth.ProviderGoogle, &fakeOAuthClient{payload: fakePayload})
			},
			input: usecase.OAuthCallbackInput{
				Provider: model.OAuthProvider("unknown"),
				Code:     "valid_code",
			},
			assert: func(t *testing.T, infra *di.Infra, got usecase.OAuthCallbackOutput, err error) {
				require.ErrorIs(t, err, usecase.ErrOAuthCallbackUnsupportedProvider)
			},
		},
		{
			name: "callback failed",
			arrange: func(t *testing.T, infra *di.Infra) *oauth.Resolver {
				return newFakeResolver(oauth.ProviderGoogle, &fakeOAuthClient{
					err: assert.AnError,
				})
			},
			input: usecase.OAuthCallbackInput{
				Provider: model.OAuthProviderGoogle,
				Code:     "invalid_code",
			},
			assert: func(t *testing.T, infra *di.Infra, got usecase.OAuthCallbackOutput, err error) {
				require.ErrorIs(t, err, usecase.ErrOAuthCallbackFailed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			infra := newInfra(t)
			resolver := tt.arrange(t, infra)

			// act
			sut := usecase.NewOAuthCallback(infra, resolver)
			got, err := sut.Do(t.Context(), tt.input)

			// assert
			tt.assert(t, infra, got, err)
		})
	}
}
