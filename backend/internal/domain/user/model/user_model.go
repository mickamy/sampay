package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"mickamy.com/sampay/internal/lib/random"
	"mickamy.com/sampay/internal/lib/ulid"
)

type User struct {
	ID        string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
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
	return nil
}
