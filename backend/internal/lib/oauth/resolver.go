package oauth

import (
	"fmt"

	"github.com/mickamy/sampay/config"
)

type Resolver struct {
	Clients map[Provider]Client
}

func NewResolverFromConfig(cfg config.OAuthConfig) *Resolver {
	return &Resolver{Clients: map[Provider]Client{
		ProviderLINE: NewLINE(cfg),
	}}
}

func (r *Resolver) Resolve(provider Provider) (Client, error) {
	client, ok := r.Clients[provider]
	if !ok {
		return nil, fmt.Errorf("oauth: unsupported provider: %s", provider)
	}
	return client, nil
}
