package model

import (
	"time"
)

type UserLinkDisplayAttribute struct {
	UserLinkID   string `gorm:"primaryKey"`
	Name         string
	DisplayOrder int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
