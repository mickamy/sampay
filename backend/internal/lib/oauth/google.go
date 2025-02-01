package oauth

import (
	"encoding/base64"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/lib/random"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type Google interface {
	AuthenticationURL() (string, error)
}

type googleClient struct {
	config *oauth2.Config
}

func NewGoogle(cfg config.OAuthConfig) Google {
	return &googleClient{
		config: &oauth2.Config{
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			Endpoint:     google.Endpoint,
			RedirectURL:  cfg.RedirectURL + "/google",
			Scopes:       []string{"openid", "profile", "email"},
		},
	}
}

func (c *googleClient) AuthenticationURL() (string, error) {
	bytes, err := random.NewBytes(16)
	if err != nil {
		return "", err
	}
	s := base64.URLEncoding.EncodeToString(bytes)
	authURL := c.config.AuthCodeURL(s, oauth2.AccessTypeOffline)
	return authURL, nil
}
