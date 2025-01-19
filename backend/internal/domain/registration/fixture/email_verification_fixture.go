package fixture

import (
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v7"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/domain/registration/model"
)

func EmailVerification(setter func(m *model.EmailVerification)) model.EmailVerification {
	m := model.EmailVerification{
		Email: gofakeit.GlobalFaker.Email(),
	}

	if setter != nil {
		setter(&m)
	}

	return m
}

func EmailVerificationRequested(setter func(m *model.EmailVerification)) model.EmailVerification {
	m := EmailVerification(func(m *model.EmailVerification) {
		if err := m.Request(config.Auth().EmailVerificationExpiresInDuration()); err != nil {
			log.Fatal(fmt.Errorf("failed to request email verification: %w", err))
		}
	})

	if setter != nil {
		setter(&m)
	}

	return m
}

func EmailVerificationVerified(setter func(m *model.EmailVerification)) model.EmailVerification {
	m := EmailVerificationRequested(func(m *model.EmailVerification) {
		if err := m.Verify(); err != nil {
			log.Fatal(fmt.Errorf("failed to verify email verification: %w", err))
		}
	})

	if setter != nil {
		setter(&m)
	}

	return m
}
