package usecase

import (
	"context"
	"errors"
	"fmt"

	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	"mickamy.com/sampay/internal/infra/storage/database"
)

var (
	ErrAuthenticateAnonymousUserSessionNotFound = errors.New("session not found")
)

type AuthenticateAnonymousUserInput struct {
	Token string
}

type AuthenticateAnonymousUserOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type AuthenticateAnonymousUser interface {
	Do(ctx context.Context, input AuthenticateAnonymousUserInput) (AuthenticateAnonymousUserOutput, error)
}

type authenticateAnonymousUser struct {
	reader                *database.Reader
	emailVerificationRepo authRepository.EmailVerification
}

func NewAuthenticateAnonymousUser(
	reader *database.Reader,
	emailVerificationRepo authRepository.EmailVerification,
) AuthenticateAnonymousUser {
	return &authenticateAnonymousUser{
		reader:                reader,
		emailVerificationRepo: emailVerificationRepo,
	}
}

func (uc *authenticateAnonymousUser) Do(ctx context.Context, input AuthenticateAnonymousUserInput) (AuthenticateAnonymousUserOutput, error) {
	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		verification, err := uc.emailVerificationRepo.WithTx(tx.ReaderDB()).FindByVerifiedToken(ctx, input.Token)
		if err != nil {
			return fmt.Errorf("failed to find email verification: %w", err)
		}
		if verification == nil {
			return errors.New("email verification not found")
		}

		return nil
	}); err != nil {
		return AuthenticateAnonymousUserOutput{}, err
	}

	return AuthenticateAnonymousUserOutput{}, nil
}
