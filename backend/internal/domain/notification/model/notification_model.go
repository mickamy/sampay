package model

import (
	"time"

	"gorm.io/gorm"

	"mickamy.com/sampay/internal/lib/ptr"
	"mickamy.com/sampay/internal/lib/ulid"
)

type Notification struct {
	ID        string
	Type      NotificationType
	UserID    string
	Subject   string
	Body      string
	CreatedAt time.Time

	ReadStatus NotificationReadStatus
}

func (m *Notification) Read() {
	m.ReadStatus = NotificationReadStatus{
		NotificationID: m.ID,
		UserID:         m.UserID,
		ReadAt:         ptr.Of(time.Now()),
	}
}

func (m *Notification) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.New()
	}
	return nil
}
