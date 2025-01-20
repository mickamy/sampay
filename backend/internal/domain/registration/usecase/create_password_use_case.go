package usecase

import (
	"context"
	"errors"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrCreatePasswordEmailVerificationInvalidToken = commonModel.NewLocalizableError(errors.New("invalid email verification token")).
							WithMessages(i18n.Config{MessageID: i18n.RegistrationUsecaseCreate_passwordErrorEmail_verification_invalid_token})
	ErrCreatePasswordEmailVerificationAlreadyConsumed = commonModel.NewLocalizableError(errors.New("email verification already consumed")).
								WithMessages(i18n.Config{MessageID: i18n.RegistrationUsecaseCreate_passwordErrorEmail_verification_already_consumed})
)

type CreatePasswordInput struct {
	Token    string
	Password string
}

type CreatePasswordOutput struct {
	Session authModel.Session
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreatePassword interface {
	Do(ctx context.Context, input CreatePasswordInput) (CreatePasswordOutput, error)
}

type createPassword struct {
	writer                *database.Writer
	emailVerificationRepo authRepository.EmailVerification
	authenticationRepo    authRepository.Authentication
	userRepo              userRepository.User
	sessionRepo           authRepository.Session
}

func NewCreatePassword(
	writer *database.Writer,
	emailVerificationRepo authRepository.EmailVerification,
	authenticationRepo authRepository.Authentication,
	userRepo userRepository.User,
	sessionRepo authRepository.Session,
) CreatePassword {
	return &createPassword{
		writer:                writer,
		emailVerificationRepo: emailVerificationRepo,
		authenticationRepo:    authenticationRepo,
		userRepo:              userRepo,
		sessionRepo:           sessionRepo,
	}
}

func (uc *createPassword) Do(ctx context.Context, input CreatePasswordInput) (CreatePasswordOutput, error) {
	var session authModel.Session
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		verification, err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).FindByVerifiedToken(ctx, input.Token, authRepository.EmailVerificationInnerJoinRequested, authRepository.EmailVerificationJoinConsumed)
		if err != nil {
			return fmt.Errorf("failed to find email verification: %w", err)
		}
		if verification == nil {
			return fmt.Errorf("email verification not found: %w", ErrCreatePasswordEmailVerificationInvalidToken)
		}
		if verification.IsConsumed() {
			return ErrCreatePasswordEmailVerificationAlreadyConsumed
		}

		if err := verification.Consume(); err != nil {
			return fmt.Errorf("failed to consume email verification: %w", err)
		}
		if err := uc.emailVerificationRepo.WithTx(tx.WriterDB()).Update(ctx, verification); err != nil {
			return fmt.Errorf("failed to update email verification: %w", err)
		}

		user := userModel.User{}
		if err := uc.userRepo.WithTx(tx.WriterDB()).Create(ctx, &user); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		auth, err := authModel.NewAuthenticationEmailPassword(user.ID, verification.Email, input.Password)
		if err != nil {
			return fmt.Errorf("failed to create authentication: %w", err)
		}

		if err := uc.authenticationRepo.WithTx(tx.WriterDB()).Create(ctx, &auth); err != nil {
			return fmt.Errorf("failed to persist user attribute: %w", err)
		}

		session, err = authModel.NewSession(user.ID)
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}

		if err := uc.sessionRepo.Create(ctx, session); err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}

		return nil
	}); err != nil {
		return CreatePasswordOutput{}, err
	}

	return CreatePasswordOutput{Session: session}, nil
}
