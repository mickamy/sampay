package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrUserAttributeDuplicated = commonModel.NewLocalizableError(
		errors.New("user attribute already exists"),
	).WithMessages(i18n.Config{MessageID: i18n.UserModelUser_attributeErrorDuplicated})
)

type UserAttribute struct {
	UserID            string `gorm:"primaryKey"`
	UsageCategoryType string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (m UserAttribute) IsZero() bool {
	return m == UserAttribute{}
}

func (m *UserAttribute) BeforeCreate(db *gorm.DB) error {
	if err := m.ValidateExistence(db); err != nil {
		return err
	}

	return nil
}

func (m UserAttribute) ValidateExistence(db *gorm.DB) error {
	var existingID string
	if err := db.Model(&UserAttribute{}).Where("user_id = ?", m.UserID).Limit(1).Pluck("user_id", &existingID).Error; err != nil {
		return fmt.Errorf("failed to validate existence: %w", err)
	}
	if existingID != "" {
		return ErrUserAttributeDuplicated
	}
	return nil
}
