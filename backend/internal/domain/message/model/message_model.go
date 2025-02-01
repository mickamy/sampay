package model

import (
	"time"

	"gorm.io/gorm"

	"mickamy.com/sampay/internal/lib/ulid"
)

type Message struct {
	ID         string
	SenderName string
	ReceiverID string
	Content    string
	CreatedAt  time.Time
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.New()
	}
	return nil
}
