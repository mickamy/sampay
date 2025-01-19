package model

import (
	"time"
)

type VerifiedEmailVerification struct {
	EmailVerificationID string `gorm:"primaryKey"`
	Token               string
	VerifiedAt          time.Time
}
