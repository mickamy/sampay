package model

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"mickamy.com/sampay/internal/lib/random"
	"mickamy.com/sampay/internal/lib/ulid"
)

var (
	ErrEmailVerificationTokenExpired = errors.New("email verification token expired")
	ErrEmailVerificationNotRequested = errors.New("email verification not requested")
	ErrEmailVerificationRequested    = errors.New("email verification requested")
	ErrEmailVerificationNotVerified  = errors.New("email verification not verified")
	ErrEmailVerificationVerified     = errors.New("email verification verified")
	ErrEmailVerificationConsumed     = errors.New("email verification consumed")
)

type EmailVerification struct {
	ID         string
	IntentType EmailVerificationIntentType
	Email      string
	CreatedAt  time.Time

	Requested *RequestedEmailVerification
	Verified  *VerifiedEmailVerification
	Consumed  *ConsumedEmailVerification
}

func (m *EmailVerification) BeforeCreate(tx *gorm.DB) error {
	if m.ID != "" {
		return nil
	}
	m.ID = ulid.New()
	return nil
}

func (m *EmailVerification) Request(expiresIn time.Duration) error {
	if m.IsConsumed() {
		return ErrEmailVerificationConsumed
	}
	if m.IsVerified() {
		return ErrEmailVerificationVerified
	}
	if m.IsRequested() {
		return nil
	}
	pin, err := random.NewPinCode(6)
	if err != nil {
		return fmt.Errorf("failed to generate pin code: %w", err)
	}
	token, err := random.NewString(32)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	now := time.Now()
	m.Requested = &RequestedEmailVerification{
		PINCode:     pin,
		Token:       token,
		RequestedAt: now,
		ExpiresAt:   now.Add(expiresIn),
	}
	return nil
}

func (m EmailVerification) IsRequested() bool {
	if m.Requested == nil {
		return false
	}
	return !m.Requested.RequestedAt.IsZero()
}

func (m *EmailVerification) Verify() error {
	if m.IsConsumed() {
		return ErrEmailVerificationConsumed
	}
	if m.IsVerified() {
		return nil
	}
	if !m.IsRequested() {
		return ErrEmailVerificationNotRequested
	}
	now := time.Now()
	if now.After(m.Requested.ExpiresAt) {
		return ErrEmailVerificationTokenExpired
	}
	token, err := random.NewString(32)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}
	m.Verified = &VerifiedEmailVerification{
		Token:      token,
		VerifiedAt: now,
	}
	return nil
}

func (m EmailVerification) IsVerified() bool {
	if m.Verified == nil {
		return false
	}
	return !m.Verified.VerifiedAt.IsZero()
}

func (m *EmailVerification) Consume() error {
	if m.IsConsumed() {
		return nil
	}
	if !m.IsRequested() {
		return ErrEmailVerificationNotRequested
	}
	if !m.IsVerified() {
		return ErrEmailVerificationNotVerified
	}
	m.Consumed = &ConsumedEmailVerification{
		ConsumedAt: time.Now(),
	}
	return nil
}

func (m EmailVerification) IsConsumed() bool {
	if m.Consumed == nil {
		return false
	}
	return !m.Consumed.ConsumedAt.IsZero()
}
