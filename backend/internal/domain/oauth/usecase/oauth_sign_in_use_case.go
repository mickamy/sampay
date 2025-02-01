package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/domain/oauth/model"
	"mickamy.com/sampay/internal/lib/oauth"
)

type OAuthSignInInput struct {
	Provider model.OAuthProvider
}

type OAuthSignInOutput struct {
	AuthenticationURL string
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type OAuthSignIn interface {
	Do(ctx context.Context, input OAuthSignInInput) (OAuthSignInOutput, error)
}

type oauthSignIn struct {
	google oauth.Google
}

func NewOAuthSignIn(
	google oauth.Google,
) OAuthSignIn {
	return &oauthSignIn{
		google: google,
	}
}

func (uc *oauthSignIn) Do(ctx context.Context, input OAuthSignInInput) (OAuthSignInOutput, error) {
	switch input.Provider {
	case model.OAuthProviderGoogle:
		url, err := uc.google.AuthenticationURL()
		if err != nil {
			return OAuthSignInOutput{}, fmt.Errorf("failed to get google authentication url: %w", err)
		}
		return OAuthSignInOutput{
			AuthenticationURL: url,
		}, nil
	default:
		return OAuthSignInOutput{}, fmt.Errorf("unsupported provider: %s", input.Provider)
	}

}
