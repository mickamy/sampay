package oauth

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"

	"github.com/mickamy/sampay/config"
)

var lineEndpoint = oauth2.Endpoint{
	AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
	TokenURL: "https://api.line.me/oauth2/v2.1/token",
}

var _ Client = (*lineClient)(nil)

type lineClient struct {
	cfg *oauth2.Config
}

func NewLINE(cfg config.OAuthConfig) Client {
	return &lineClient{
		cfg: &oauth2.Config{
			ClientID:     cfg.LINEChannelID,
			ClientSecret: cfg.LINEChannelSecret,
			Endpoint:     lineEndpoint,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"profile", "openid", "email"},
		},
	}
}

func (c *lineClient) AuthenticationURL() (string, error) {
	s, err := generateState(ProviderLINE)
	if err != nil {
		return "", fmt.Errorf("oauth: failed to generate state: %w", err)
	}
	authURL := c.cfg.AuthCodeURL(s, oauth2.AccessTypeOffline)
	return authURL, nil
}

func (c *lineClient) Callback(ctx context.Context, code string) (Payload, error) {
	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return Payload{}, fmt.Errorf("oauth: failed to exchange token: %w", err)
	}

	httpClient := c.cfg.Client(ctx, token)
	resp, err := httpClient.Get("https://api.line.me/v2/profile")
	if err != nil {
		return Payload{}, fmt.Errorf("oauth: failed to get profile: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var u struct {
		UserID      string `json:"userId"`
		DisplayName string `json:"displayName"`
		PictureURL  string `json:"pictureUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return Payload{}, fmt.Errorf("oauth: failed to decode profile: %w", err)
	}

	return Payload{
		Provider: ProviderLINE,
		UID:      u.UserID,
		Name:     u.DisplayName,
		Picture:  u.PictureURL,
	}, nil
}
