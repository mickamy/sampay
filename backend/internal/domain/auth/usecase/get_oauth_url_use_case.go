package usecase

import (
	"context"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/lib/oauth"
)

type GetOAuthURLInput struct {
	Provider model.OAuthProvider
}

type GetOAuthURLOutput struct {
	AuthenticationURL string
}

type GetOAuthURL interface {
	Do(ctx context.Context, input GetOAuthURLInput) (GetOAuthURLOutput, error)
}

type getOAuthURL struct {
	_        GetOAuthURL     `inject:"returns"`
	resolver *oauth.Resolver `inject:"param"`
}

func (uc *getOAuthURL) Do(ctx context.Context, input GetOAuthURLInput) (GetOAuthURLOutput, error) {
	client, err := uc.resolver.Resolve(oauth.Provider(input.Provider))
	if err != nil {
		return GetOAuthURLOutput{}, errx.
			Wrap(err, "failed to resolve oauth client").
			With("provider", input.Provider).
			WithCode(errx.Internal)
	}

	url, err := client.AuthenticationURL()
	if err != nil {
		return GetOAuthURLOutput{}, errx.
			Wrap(err, "failed to get authentication url").
			With("provider", input.Provider).
			WithCode(errx.Internal)
	}

	return GetOAuthURLOutput{AuthenticationURL: url}, nil
}
