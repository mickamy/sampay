package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	oauthModel "mickamy.com/sampay/internal/domain/oauth/model"
	"mickamy.com/sampay/internal/lib/oauth"
	"mickamy.com/sampay/internal/lib/passwd"
	"mickamy.com/sampay/internal/lib/ulid"
)

type AuthenticationType string

const (
	AuthenticationTypePassword AuthenticationType = "password"
	AuthenticationTypeGoogle   AuthenticationType = "google"
)

var (
	ErrAuthenticationDifferentType = fmt.Errorf("different authentication type")
)

type Authentication struct {
	ID         string
	UserID     string
	Type       AuthenticationType
	Identifier string
	Secret     string
	MFAEnabled bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewAuthenticationEmailPassword(userID, email, password string) (Authentication, error) {
	hash, err := passwd.New(password, 16)
	if err != nil {
		return Authentication{}, fmt.Errorf("failed to hash password: %w", err)
	}
	return Authentication{
		UserID:     userID,
		Type:       AuthenticationTypePassword,
		Identifier: email,
		Secret:     hash,
	}, nil
}

func NewAuthenticationOAuth(payload oauth.Payload) (Authentication, error) {
	var provider AuthenticationType
	switch payload.Provider.String() {
	case oauthModel.OAuthProviderGoogle.String():
		provider = AuthenticationTypeGoogle
	default:
		return Authentication{}, fmt.Errorf("unsupported provider: %s", payload.Provider)
	}
	return Authentication{
		Type:       provider,
		Identifier: payload.UID,
	}, nil
}

func (m Authentication) AuthenticateByEmailAndPassword(email string, password string) (bool, error) {
	if m.Type != AuthenticationTypePassword {
		return false, ErrAuthenticationDifferentType
	}
	if m.Identifier != email {
		return false, nil
	}
	return passwd.Verify(password, m.Secret)
}

func (m *Authentication) ResetPassword(password string) error {
	hash, err := passwd.New(password, 16)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	m.Secret = hash
	return nil
}

func (m *Authentication) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.New()
	}
	return nil
}

type Authentications []Authentication

func (ms Authentications) FindByType(t AuthenticationType) *Authentication {
	for i := range ms {
		if ms[i].Type == t {
			return &ms[i]
		}
	}
	return nil
}
