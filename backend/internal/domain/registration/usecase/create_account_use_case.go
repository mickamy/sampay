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
	ErrCreateAccountEmailAlreadyExists = commonModel.
		NewLocalizableError(errors.New("email already exists")).
		WithMessages(i18n.Config{MessageID: i18n.RegistrationUsecaseCreate_accountErrorEmail_already_exists})
)

type CreateAccountInput struct {
	Email    string
	Password string
}

type CreateAccountOutput struct {
	authModel.Session
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreateAccount interface {
	Do(ctx context.Context, input CreateAccountInput) (CreateAccountOutput, error)
}

type createAccount struct {
	writer             *database.Writer
	authenticationRepo authRepository.Authentication
	sessionRepo        authRepository.Session
	userRepo           userRepository.User
}

func NewCreateAccount(
	writer *database.Writer,
	authenticationRepo authRepository.Authentication,
	sessionRepo authRepository.Session,
	userRepo userRepository.User,
) CreateAccount {
	return &createAccount{
		writer:             writer,
		authenticationRepo: authenticationRepo,
		sessionRepo:        sessionRepo,
		userRepo:           userRepo,
	}
}

func (uc *createAccount) Do(ctx context.Context, input CreateAccountInput) (CreateAccountOutput, error) {
	var session authModel.Session

	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		existingAuth, err := uc.authenticationRepo.WithTx(tx.WriterDB()).FindByTypeAndIdentifier(ctx, authModel.AuthenticationTypePassword, input.Email)
		if err != nil {
			return fmt.Errorf("failed to find authentication: %w", err)
		}
		if existingAuth != nil {
			return ErrCreateAccountEmailAlreadyExists
		}

		user := userModel.User{}
		if err := uc.userRepo.WithTx(tx.WriterDB()).Create(ctx, &user); err != nil {
			return fmt.Errorf("failed to persist user: %w", err)
		}

		auth, err := authModel.NewAuthenticationEmailPassword(user.ID, input.Email, input.Password)
		if err != nil {
			return fmt.Errorf("failed to create authentication: %w", err)
		}

		if err := uc.authenticationRepo.WithTx(tx.WriterDB()).Create(ctx, &auth); err != nil {
			return fmt.Errorf("failed to persist authentication: %w", err)
		}

		session, err = authModel.NewSession(user.ID)
		if err != nil {
			return fmt.Errorf("failed to create session: %w", err)
		}

		if err := uc.sessionRepo.Create(ctx, session); err != nil {
			return fmt.Errorf("failed to persist session: %w", err)
		}

		return nil
	}); err != nil {
		return CreateAccountOutput{}, err
	}

	return CreateAccountOutput{Session: session}, nil
}
