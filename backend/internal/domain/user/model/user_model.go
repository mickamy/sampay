package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/lib/random"
	"mickamy.com/sampay/internal/lib/ulid"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrUserSlugAlreadyTaken = commonModel.NewLocalizableError(errors.New("slug already taken")).
				WithMessages(i18n.Config{MessageID: i18n.UserModelUserErrorSlug_already_taken})
	ErrUserEmailAlreadyTaken = commonModel.NewLocalizableError(errors.New("email already taken")).
					WithMessages(i18n.Config{MessageID: i18n.UserModelUserErrorEmail_already_taken})
)

type User struct {
	ID        string
	Slug      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time

	Attribute UserAttribute
	Profile   UserProfile
	Links     []UserLink
}

func (m User) validateSlugExistence(tx *gorm.DB) error {
	var existingID string
	if err := tx.Model(User{}).Where("slug = ?", m.Slug).Limit(1).Pluck("id", &existingID).Error; err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if existingID != "" && existingID != m.ID {
		return ErrUserSlugAlreadyTaken
	}
	return nil
}

func (m User) validateEmailExistence(tx *gorm.DB) error {
	var existingID string
	if err := tx.Model(User{}).Where("email = ?", m.Email).Limit(1).Pluck("id", &existingID).Error; err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if existingID != "" && existingID != m.ID {
		return ErrUserEmailAlreadyTaken
	}
	return nil
}

func (m *User) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.New()
	}
	if m.Slug == "" {
		var err error
		m.Slug, err = random.NewString(16)
		if err != nil {
			return fmt.Errorf("failed to generate slug: %w", err)
		}
	}

	if err := m.validateSlugExistence(tx); err != nil {
		return err
	}
	if err := m.validateEmailExistence(tx); err != nil {
		return err
	}

	return nil
}

func (m *User) BeforeUpdate(tx *gorm.DB) error {
	if err := m.validateSlugExistence(tx); err != nil {
		return err
	}
	if err := m.validateEmailExistence(tx); err != nil {
		return err
	}

	return nil
}
