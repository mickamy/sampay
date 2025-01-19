package usecase

import (
	"context"
	"errors"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrCreatePasswordEmailVerificationAlreadyConsumed = commonModel.NewLocalizableError(errors.New("email verification already consumed")).
		WithMessages(i18n.Config{MessageID: i18n.RegistrationUsecaseCreate_passwordErrorEmail_verification_already_consumed})
)

type CreatePasswordInput struct {
	Email    string
	Password string
}

type CreatePasswordOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreatePassword interface {
	Do(ctx context.Context, input CreatePasswordInput) (CreatePasswordOutput, error)
}

type createPassword struct {
	writer                *database.Writer
	emailVerificationRepo authRepository.EmailVerification
	authenticationRepo    authRepository.Authentication
}

func NewCreatePassword(
	writer *database.Writer,
	emailVerificationRepo authRepository.EmailVerification,
	authenticationRepo authRepository.Authentication,
) CreatePassword {
	return &createPassword{
		writer:                writer,
		emailVerificationRepo: emailVerificationRepo,
		authenticationRepo:    authenticationRepo,
	}
}

func (uc *createPassword) Do(ctx context.Context, input CreatePasswordInput) (CreatePasswordOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		verification, err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).FindByEmail(ctx, input.Email, authRepository.EmailVerificationJoinVerified, authRepository.EmailVerificationJoinConsumed)
		if err != nil {
			return fmt.Errorf("failed to find email verification: %w", err)
		}
		if verification == nil {
			return fmt.Errorf("email verification not found: %w", ErrCreatePasswordEmailVerificationAlreadyConsumed)
		}
		if verification.IsConsumed() {
			return ErrCreatePasswordEmailVerificationAlreadyConsumed
		}

		m, err := authModel.NewAuthenticationEmailPassword(contexts.MustAuthenticatedUserID(ctx), input.Email, input.Password)
		if err != nil {
			return fmt.Errorf("failed to create authentication model: %w", err)
		}

		if err := uc.authenticationRepo.WithTx(tx.WriterDB()).Create(ctx, &m); err != nil {
			return fmt.Errorf("failed to persist user attribute: %w", err)
		}

		return nil
	}); err != nil {
		return CreatePasswordOutput{}, err
	}

	return CreatePasswordOutput{}, nil
}
