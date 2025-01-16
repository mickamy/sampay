package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/mickamy/go-sqs-worker/producer"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	registrationRepository "mickamy.com/sampay/internal/domain/registration/repository"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrVerifyEmailInvalidToken = commonModel.
		NewLocalizableError(errors.New("invalid token")).
		WithMessages(i18n.Config{MessageID: i18n.RegistrationUsecaseVerify_emailErrorInvalid_token})
)

type VerifyEmailInput struct {
	Token string
}

type VerifyEmailOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type VerifyEmail interface {
	Do(ctx context.Context, input VerifyEmailInput) (VerifyEmailOutput, error)
}

type verifyEmail struct {
	writer                *database.Writer
	producer              *producer.Producer
	authenticationRepo    authRepository.Authentication
	emailVerificationRepo registrationRepository.EmailVerification
}

func NewVerifyEmail(
	writer *database.Writer,
	producer *producer.Producer,
	authenticationRepo authRepository.Authentication,
	emailVerificationRepo registrationRepository.EmailVerification,
) VerifyEmail {
	return &verifyEmail{
		writer:                writer,
		producer:              producer,
		authenticationRepo:    authenticationRepo,
		emailVerificationRepo: emailVerificationRepo,
	}
}

func (uc *verifyEmail) Do(ctx context.Context, input VerifyEmailInput) (VerifyEmailOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		var err error
		verification, err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).FindByToken(
			ctx,
			input.Token,
			registrationRepository.EmailVerificationJoinRequested,
			registrationRepository.EmailVerificationPreloadVerified,
			registrationRepository.EmailVerificationNotConsumed,
		)
		if err != nil {
			return fmt.Errorf("failed to find email verification: %w", err)
		}
		if verification == nil {
			return ErrVerifyEmailInvalidToken
		}

		if verification.IsVerified() {
			return nil
		}
		if err := verification.Verify(); err != nil {
			return fmt.Errorf("failed to verify email verification: %w", err)
		}
		if err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).Update(ctx, verification); err != nil {
			return fmt.Errorf("failed to update email verification: %w", err)
		}
		return nil
	}); err != nil {
		return VerifyEmailOutput{}, err
	}

	return VerifyEmailOutput{}, nil
}
