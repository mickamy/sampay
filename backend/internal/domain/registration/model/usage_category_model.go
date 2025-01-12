package model

import (
	"time"
)

type UsageCategory struct {
	Type         string `gorm:"primaryKey"`
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
