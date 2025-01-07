package model

import (
	"time"
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
