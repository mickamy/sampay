package model

import (
	"time"
)

type NotificationReadStatus struct {
	NotificationID string `gorm:"primaryKey"`
	UserID         string `gorm:"primaryKey"`
	ReadAt         *time.Time
}
