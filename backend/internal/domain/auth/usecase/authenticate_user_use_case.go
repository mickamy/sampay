package usecase

import (
	"context"
	"errors"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/jwt"
)

var (
	ErrAuthenticateUserSessionNotFound = errors.New("session not found")
	ErrAuthenticateUserUserNotFound    = errors.New("user not found")
)

type AuthenticateUserInput struct {
	AccessToken string
}

type AuthenticateUserOutput struct {
	User userModel.User
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type AuthenticateUser interface {
	Do(ctx context.Context, input AuthenticateUserInput) (AuthenticateUserOutput, error)
}

type authenticateUser struct {
	reader      *database.Reader
	sessionRepo authRepository.Session
	userRepo    userRepository.User
}

func NewAuthenticateUser(
	reader *database.Reader,
	sessionRepo authRepository.Session,
	userRepo userRepository.User,
) AuthenticateUser {
	return &authenticateUser{
		reader:      reader,
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

func (uc *authenticateUser) Do(ctx context.Context, input AuthenticateUserInput) (AuthenticateUserOutput, error) {
	userID, err := jwt.ExtractID(input.AccessToken)
	if err != nil {
		return AuthenticateUserOutput{}, fmt.Errorf("failed to extract user id from access token: %w", err)
	}

	exists, err := uc.sessionRepo.AccessTokenExists(ctx, userID, input.AccessToken)
	if err != nil {
		return AuthenticateUserOutput{}, fmt.Errorf("failed to check access token existence: %w", err)
	}
	if !exists {
		return AuthenticateUserOutput{}, ErrAuthenticateUserSessionNotFound
	}

	var user userModel.User
	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		u, err := uc.userRepo.WithTx(tx.ReaderDB()).FindByID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to find user: %w", err)
		}
		if u == nil {
			return ErrAuthenticateUserUserNotFound
		}

		user = *u

		return nil
	}); err != nil {
		return AuthenticateUserOutput{}, err
	}

	return AuthenticateUserOutput{User: user}, nil
}
