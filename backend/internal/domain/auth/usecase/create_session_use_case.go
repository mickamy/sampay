package usecase

import (
	"context"
	"errors"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrCreateSessionPasswordNotMatch = commonModel.
		NewLocalizableError(errors.New("password not match")).
		WithMessages(i18n.Config{MessageID: i18n.AuthUsecaseCreate_sessionInvalid_email_password})
)

type CreateSessionInput struct {
	Email    string
	Password string
}

type CreateSessionOutput struct {
	Session authModel.Session
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreateSession interface {
	Do(ctx context.Context, input CreateSessionInput) (CreateSessionOutput, error)
}

type createSession struct {
	reader             *database.Reader
	authenticationRepo authRepository.Authentication
	sessionRepo        authRepository.Session
	userRepo           userRepository.User
}

func NewCreateSession(
	reader *database.Reader,
	authenticationRepo authRepository.Authentication,
	sessionRepo authRepository.Session,
	userRepo userRepository.User,
) CreateSession {
	return &createSession{
		reader:             reader,
		authenticationRepo: authenticationRepo,
		sessionRepo:        sessionRepo,
		userRepo:           userRepo,
	}
}

func (uc *createSession) Do(ctx context.Context, input CreateSessionInput) (CreateSessionOutput, error) {
	var userID string

	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		user, err := uc.userRepo.WithTx(tx.ReaderDB()).FindByEmail(ctx, input.Email)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}
		if user == nil {
			return ErrCreateSessionPasswordNotMatch
		}

		auths, err := uc.authenticationRepo.WithTx(tx.ReaderDB()).ListByUserID(ctx, user.ID)
		if err != nil {
			return err
		}

		emailPasswordAuth := authModel.Authentications(auths).FindByType(authModel.AuthenticationTypeEmailPassword)
		if emailPasswordAuth == nil {
			return ErrCreateSessionPasswordNotMatch
		}

		matched, err := emailPasswordAuth.AuthenticateByEmailAndPassword(input.Email, input.Password)
		if err != nil {
			return err
		}
		if !matched {
			return ErrCreateSessionPasswordNotMatch
		}

		userID = user.ID

		return nil
	}); err != nil {
		return CreateSessionOutput{}, err
	}

	session, err := authModel.NewSession(userID)
	if err != nil {
		return CreateSessionOutput{}, fmt.Errorf("failed to create session: %w", err)
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return CreateSessionOutput{}, fmt.Errorf("failed to persist session: %w", err)
	}

	return CreateSessionOutput{Session: session}, nil
}
