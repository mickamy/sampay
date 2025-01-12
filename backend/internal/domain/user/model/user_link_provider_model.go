package model

import (
	"time"
)

type UserLinkProvider struct {
	Type         UserLinkProviderType `gorm:"primaryKey"`
	DisplayOrder int

	CreatedAt time.Time
	UpdatedAt time.Time
}
