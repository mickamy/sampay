package model

import (
	"time"
)

type RequestedEmailVerification struct {
	EmailVerificationID string `gorm:"primaryKey"`
	Token               string
	RequestedAt         time.Time
	ExpiresAt           time.Time
}
