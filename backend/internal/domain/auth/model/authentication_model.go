package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"mickamy.com/sampay/internal/lib/passwd"
	"mickamy.com/sampay/internal/lib/ulid"
)

type AuthenticationType string

const (
	AuthenticationTypeEmailPassword AuthenticationType = "email_password"
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
		Type:       AuthenticationTypeEmailPassword,
		Identifier: email,
		Secret:     hash,
	}, nil
}

func (m Authentication) AuthenticateByEmailAndPassword(email string, password string) (bool, error) {
	if m.Type != AuthenticationTypeEmailPassword {
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
