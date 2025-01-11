package usecase

import (
	"context"
	"errors"
	"fmt"

	authModel "mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/auth/repository"
	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/lib/jwt"
	"mickamy.com/sampay/internal/misc/i18n"
)

var (
	ErrDeleteSessionNotFound = commonModel.
					NewLocalizableError(errors.New("session not found")).
					WithMessages(i18n.Config{MessageID: "auth.usecase.error.invalid_access_refresh_token"})
	ErrDeleteSessionTokenMismatch = commonModel.
					NewLocalizableError(errors.New("token mismatch")).
					WithMessages(i18n.Config{MessageID: "auth.usecase.error.invalid_access_refresh_token"})
	ErrDeleteSessionDeletingTokensFailed = errors.New("deleting tokens failed")
)

type DeleteSessionInput struct {
	AccessToken  string
	RefreshToken string
}

type DeleteSessionOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type DeleteSession interface {
	Do(ctx context.Context, input DeleteSessionInput) (DeleteSessionOutput, error)
}

type deleteSession struct {
	sessionRepo repository.Session
}

func NewDeleteSession(
	sessionRepo repository.Session,
) DeleteSession {
	return &deleteSession{
		sessionRepo: sessionRepo,
	}
}

func (uc *deleteSession) Do(ctx context.Context, input DeleteSessionInput) (DeleteSessionOutput, error) {
	userID, err := jwt.ExtractID(input.AccessToken)
	if err != nil {
		return DeleteSessionOutput{}, errors.Join(ErrDeleteSessionNotFound, fmt.Errorf("failed to extract user id from access token: %w", err))
	}

	userIDFromRefresh, err := jwt.ExtractID(input.RefreshToken)
	if err != nil {
		return DeleteSessionOutput{}, errors.Join(ErrDeleteSessionNotFound, fmt.Errorf("failed to extract user id from refresh token: %w", err))
	}

	if userID != userIDFromRefresh {
		return DeleteSessionOutput{}, errors.Join(ErrDeleteSessionTokenMismatch)
	}

	if err := uc.sessionRepo.Delete(ctx, authModel.Session{
		UserID: userID,
		Tokens: jwt.Tokens{
			Access: jwt.Token{
				Value: input.AccessToken,
			},
			Refresh: jwt.Token{
				Value: input.RefreshToken,
			},
		},
	}); err != nil {
		return DeleteSessionOutput{}, errors.Join(ErrDeleteSessionDeletingTokensFailed, err)
	}
	return DeleteSessionOutput{}, nil
}
