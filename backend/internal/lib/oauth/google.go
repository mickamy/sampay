package oauth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/lib/random"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type Google interface {
	AuthenticationURL() (string, error)
	Validate(ctx context.Context, code string) (*Payload, error)
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

func (c *googleClient) Validate(ctx context.Context, code string) (*Payload, error) {
	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to get google token: %w", err)
	}
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, fmt.Errorf("id_token not found in google token: %v", token)
	}

	payload, err := idtoken.Validate(ctx, idToken, c.config.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to validate google token: %w", err)
	}

	claims := payload.Claims
	sub, ok1 := claims["sub"].(string)
	name, ok2 := claims["name"].(string)
	email, ok3 := claims["email"].(string)
	picture, ok4 := claims["picture"].(string)

	if !ok1 || sub == "" {
		return nil, errors.New("google token missing 'sub' claim")
	}
	if !ok2 || name == "" {
		return nil, errors.New("google token missing 'name' claim")
	}
	if !ok3 || email == "" {
		return nil, errors.New("google token missing 'email' claim")
	}
	if !ok4 || picture == "" {
		// picture is optional
	}

	return &Payload{
		Provider: ProviderGoogle,
		UID:      sub,
		Name:     name,
		Email:    email,
		Picture:  picture,
	}, nil
}
