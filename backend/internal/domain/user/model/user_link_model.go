package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"mickamy.com/sampay/internal/lib/ulid"
)

var (
	ErrUserLinkInvalidURI = errors.New("invalid uri")
)

type UserLink struct {
	ID           string
	UserID       string
	ProviderType UserLinkProviderType
	URI          string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	DisplayAttribute UserLinkDisplayAttribute
}

func (m UserLink) Validate() error {
	if !m.ProviderType.MatchString(m.URI) {
		return errors.Join(ErrUserLinkInvalidURI, fmt.Errorf("uri=[%s]", m.URI))
	}
	return nil
}

func (m *UserLink) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.New()
	}
	return nil
}
