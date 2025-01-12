package model

import (
	"time"

	commonModel "mickamy.com/sampay/internal/domain/common/model"
)

type UserProfile struct {
	UserID    string `gorm:"primaryKey"`
	Name      string
	Bio       *string
	ImageID   *string
	CreatedAt time.Time
	UpdatedAt time.Time

	Image *commonModel.S3Object
}

func (m UserProfile) IsZero() bool {
	return m == UserProfile{}
}
