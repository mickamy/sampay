package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/mickamy/sampay/config"
)

var _ Client = (*googleClient)(nil)

type googleClient struct {
	cfg *oauth2.Config
}

func NewGoogle(cfg config.OAuthConfig) Client {
	return &googleClient{
		cfg: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"openid", "email", "profile"},
		},
	}
}

func (c *googleClient) AuthenticationURL() (string, error) {
	s, err := generateState(ProviderGoogle)
	if err != nil {
		return "", fmt.Errorf("oauth: failed to generate state: %w", err)
	}
	authURL := c.cfg.AuthCodeURL(s, oauth2.AccessTypeOffline)
	return authURL, nil
}

func (c *googleClient) Callback(ctx context.Context, code string) (Payload, error) {
	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return Payload{}, fmt.Errorf("oauth: failed to exchange token: %w", err)
	}

	httpClient := c.cfg.Client(ctx, token)
	resp, err := httpClient.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return Payload{}, fmt.Errorf("oauth: failed to get user info: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return Payload{}, fmt.Errorf("oauth: google userinfo returned status %d", resp.StatusCode)
	}

	var u struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return Payload{}, fmt.Errorf("oauth: failed to decode user info: %w", err)
	}

	return Payload{
		Provider: ProviderGoogle,
		UID:      u.ID,
		Name:     u.Name,
		Email:    u.Email,
		Picture:  u.Picture,
	}, nil
}
