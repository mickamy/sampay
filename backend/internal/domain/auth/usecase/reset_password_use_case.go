package usecase

import (
	"context"
	"errors"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrResetPasswordEmailVerificationInvalidToken = commonModel.NewLocalizableError(errors.New("invalid email verification token")).
							WithMessages(i18n.Config{MessageID: i18n.AuthUsecaseReset_passwordErrorEmail_verification_invalid_token})
	ErrResetPasswordEmailVerificationAlreadyConsumed = commonModel.NewLocalizableError(errors.New("email verification already consumed")).
								WithMessages(i18n.Config{MessageID: i18n.AuthUsecaseReset_passwordErrorEmail_verification_already_consumed})
)

type ResetPasswordInput struct {
	Token    string
	Password string
}

type ResetPasswordOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type ResetPassword interface {
	Do(ctx context.Context, input ResetPasswordInput) (ResetPasswordOutput, error)
}

type resetPassword struct {
	writer                *database.Writer
	emailVerificationRepo authRepository.EmailVerification
	authenticationRepo    authRepository.Authentication
}

func NewResetPassword(
	writer *database.Writer,
	emailVerificationRepo authRepository.EmailVerification,
	authenticationRepo authRepository.Authentication,
) ResetPassword {
	return &resetPassword{
		writer:                writer,
		emailVerificationRepo: emailVerificationRepo,
		authenticationRepo:    authenticationRepo,
	}
}

func (uc *resetPassword) Do(ctx context.Context, input ResetPasswordInput) (ResetPasswordOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		verification, err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).FindByVerifiedToken(ctx, input.Token, authRepository.EmailVerificationInnerJoinRequested, authRepository.EmailVerificationJoinConsumed)
		if err != nil {
			return fmt.Errorf("failed to find email verification: %w", err)
		}
		if verification == nil {
			return fmt.Errorf("email verification not found: %w", ErrResetPasswordEmailVerificationInvalidToken)
		}
		if verification.IsConsumed() {
			return ErrResetPasswordEmailVerificationAlreadyConsumed
		}

		if err := verification.Consume(); err != nil {
			return fmt.Errorf("failed to consume email verification: %w", err)
		}
		if err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).Update(ctx, verification); err != nil {
			return fmt.Errorf("failed to update email verification: %w", err)
		}

		auth, err := uc.authenticationRepo.WithTx(tx.WriterDB()).FindByTypeAndIdentifier(ctx, authModel.AuthenticationTypePassword, verification.Email)
		if err != nil {
			return fmt.Errorf("failed to find authentication: %w", err)
		}
		if auth == nil {
			return errors.New("authentication not found")
		}

		if err := auth.ResetPassword(input.Password); err != nil {
			return fmt.Errorf("failed to reset password: %w", err)
		}
		if err := uc.authenticationRepo.WithTx(tx.WriterDB()).Update(ctx, auth); err != nil {
			return fmt.Errorf("failed to update authentication: %w", err)
		}

		return nil
	}); err != nil {
		return ResetPasswordOutput{}, err
	}

	return ResetPasswordOutput{}, nil
}
