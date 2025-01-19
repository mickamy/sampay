package usecase

import (
	"context"
	"errors"
	"fmt"

	"mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/lib/jwt"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrRefreshSessionTokenNotFound = commonModel.
		NewLocalizableError(errors.New("token not found")).
		WithMessages(i18n.Config{
			MessageID: i18n.AuthUsecaseRefresh_sessionInvalid_refresh_token,
		})
)

type RefreshSessionInput struct {
	RefreshToken string
}

type RefreshSessionOutput struct {
	Tokens jwt.Tokens
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type RefreshSession interface {
	Do(ctx context.Context, input RefreshSessionInput) (RefreshSessionOutput, error)
}

type refreshSession struct {
	sessionRepo repository.Session
}

func NewRefreshSession(
	sessionRepo repository.Session,
) RefreshSession {
	return &refreshSession{
		sessionRepo: sessionRepo,
	}
}

func (uc *refreshSession) Do(ctx context.Context, input RefreshSessionInput) (RefreshSessionOutput, error) {
	id, err := jwt.ExtractID(input.RefreshToken)
	if err != nil {
		return RefreshSessionOutput{}, errors.Join(ErrRefreshSessionTokenNotFound, fmt.Errorf("failed to extract user ID from refresh token: %w", err))
	}

	exists, err := uc.sessionRepo.RefreshTokenExists(ctx, id, input.RefreshToken)
	if err != nil {
		return RefreshSessionOutput{}, errors.Join(ErrRefreshSessionTokenNotFound, fmt.Errorf("failed to check if refresh token exists: %w", err))
	}
	if !exists {
		return RefreshSessionOutput{}, ErrRefreshSessionTokenNotFound
	}

	session, err := model.NewSession(id)
	if err != nil {
		return RefreshSessionOutput{}, fmt.Errorf("failed to create new session: %w", err)
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return RefreshSessionOutput{}, fmt.Errorf("failed to persist session: %w", err)
	}

	return RefreshSessionOutput{Tokens: session.Tokens}, nil
}
