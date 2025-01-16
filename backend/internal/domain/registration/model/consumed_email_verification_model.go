package model

import (
	"time"
)

type ConsumedEmailVerification struct {
	EmailVerificationID string `gorm:"primaryKey"`
	ConsumedAt          time.Time
}
