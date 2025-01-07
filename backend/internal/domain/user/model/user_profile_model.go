package model

import (
	"time"
)

type UserProfile struct {
	UserID    string `gorm:"primaryKey"`
	Name      string
	Bio       *string
	ImageID   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m UserProfile) IsZero() bool {
	return m == UserProfile{}
}
