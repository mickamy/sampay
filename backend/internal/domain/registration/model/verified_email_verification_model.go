package model

import (
	"time"
)

type VerifiedEmailVerification struct {
	EmailVerificationID string `gorm:"primaryKey"`
	VerifiedAt          time.Time
}
