package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	registrationModel "mickamy.com/sampay/internal/domain/registration/model"
	registrationRepository "mickamy.com/sampay/internal/domain/registration/repository"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrRequestEmailVerificationEmailAlreadyExists = commonModel.
		NewLocalizableError(errors.New("email already exists")).
		WithMessages(i18n.Config{MessageID: i18n.RegistrationUsecaseCommonErrorEmail_already_exists})
)

type RequestEmailVerificationInput struct {
	Email    string
	Password string
}

type RequestEmailVerificationOutput struct {
	Token     string
	ExpiresAt time.Time
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type RequestEmailVerification interface {
	Do(ctx context.Context, input RequestEmailVerificationInput) (RequestEmailVerificationOutput, error)
}

type requestEmailVerification struct {
	writer                *database.Writer
	authenticationRepo    authRepository.Authentication
	emailVerificationRepo registrationRepository.EmailVerification
}

func NewRequestEmailVerification(
	writer *database.Writer,
	authenticationRepo authRepository.Authentication,
	emailVerificationRepo registrationRepository.EmailVerification,
) RequestEmailVerification {
	return &requestEmailVerification{
		writer:                writer,
		authenticationRepo:    authenticationRepo,
		emailVerificationRepo: emailVerificationRepo,
	}
}

func (uc *requestEmailVerification) Do(ctx context.Context, input RequestEmailVerificationInput) (RequestEmailVerificationOutput, error) {
	m := registrationModel.EmailVerification{Email: input.Email}
	if err := m.Request(config.Auth().EmailVerificationExpiresInDuration()); err != nil {
		return RequestEmailVerificationOutput{}, fmt.Errorf("failed to request email verification: %w", err)
	}
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		exists, err := uc.authenticationRepo.WithTx(tx.WriterDB()).ExistsByTypeAndIdentifier(ctx, authModel.AuthenticationTypeEmailPassword, m.Email)
		if err != nil {
			return fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return errors.Join(ErrRequestEmailVerificationEmailAlreadyExists, fmt.Errorf("email=[%s]", m.Email))
		}

		if err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).Create(ctx, &m); err != nil {
			return fmt.Errorf("failed to create email verification: %w", err)
		}

		// TODO: send email asynchronously

		return nil
	}); err != nil {
		return RequestEmailVerificationOutput{}, err
	}
	return RequestEmailVerificationOutput{Token: m.Requested.Token, ExpiresAt: m.Requested.ExpiresAt}, nil
}
