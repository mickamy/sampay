package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/lib/logger"
)

var lineEndpoint = oauth2.Endpoint{ //nolint:gosec // not credentials, just API endpoint URLs
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
	authURL := c.cfg.AuthCodeURL(s, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("bot_prompt", "aggressive"))
	return authURL, nil
}

func (c *lineClient) Callback(ctx context.Context, code string) (Payload, error) {
	token, err := c.cfg.Exchange(ctx, code)
	if err != nil {
		return Payload{}, fmt.Errorf("oauth: failed to exchange token: %w", err)
	}

	httpClient := c.cfg.Client(ctx, token)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.line.me/v2/profile", nil)
	if err != nil {
		return Payload{}, fmt.Errorf("oauth: failed to create request: %w", err)
	}
	resp, err := httpClient.Do(req) //nolint:gosec // URL is a constant, not user input
	if err != nil {
		return Payload{}, fmt.Errorf("oauth: failed to get profile: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return Payload{}, fmt.Errorf("oauth: line profile returned status %d", resp.StatusCode)
	}

	var u struct {
		UserID      string `json:"userId"`      //nolint:tagliatelle // LINE API response format
		DisplayName string `json:"displayName"` //nolint:tagliatelle // LINE API response format
		PictureURL  string `json:"pictureUrl"`  //nolint:tagliatelle // LINE API response format
	}
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return Payload{}, fmt.Errorf("oauth: failed to decode profile: %w", err)
	}

	lineFriend := c.fetchFriendship(ctx, httpClient)

	return Payload{
		Provider:   ProviderLINE,
		UID:        u.UserID,
		Name:       u.DisplayName,
		Picture:    u.PictureURL,
		LineFriend: lineFriend,
	}, nil
}

func (c *lineClient) fetchFriendship(ctx context.Context, httpClient *http.Client) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.line.me/friendship/v1/status", nil)
	if err != nil {
		logger.Warn(ctx, "failed to create friendship request", "err", err)
		return false
	}
	resp, err := httpClient.Do(req) //nolint:gosec // URL is a constant, not user input
	if err != nil {
		logger.Warn(ctx, "failed to fetch friendship status", "err", err)
		return false
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		logger.Warn(ctx, "friendship API returned non-OK status", "status", resp.StatusCode)
		return false
	}

	var f struct {
		FriendFlag bool `json:"friendFlag"` //nolint:tagliatelle // LINE API response format
	}
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		logger.Warn(ctx, "failed to decode friendship response", "err", err)
		return false
	}
	return f.FriendFlag
}
