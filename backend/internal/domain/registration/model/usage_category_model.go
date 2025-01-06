package model

import (
	"time"
)

type UsageCategory struct {
	CategoryType string `gorm:"primaryKey"`
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
