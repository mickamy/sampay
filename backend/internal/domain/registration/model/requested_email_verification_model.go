package model

import (
	"time"

	"mickamy.com/sampay/config"
)

type RequestedEmailVerification struct {
	EmailVerificationID string `gorm:"primaryKey"`
	Token               string
	RequestedAt         time.Time
	ExpiresAt           time.Time
}

func (m RequestedEmailVerification) URL() string {
	return config.WEB().BaseURL + "/auth/verify-email?token=" + m.Token
}
