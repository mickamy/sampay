package model

import (
	"time"
)

type RequestedEmailVerification struct {
	EmailVerificationID string `gorm:"primaryKey"`
	PINCode             string
	RequestedAt         time.Time
	ExpiresAt           time.Time
}
