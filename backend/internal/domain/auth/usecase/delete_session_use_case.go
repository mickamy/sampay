package usecase

import (
	"context"
	"errors"
	"fmt"

	authModel "mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/auth/repository"
	"mickamy.com/sampay/internal/lib/jwt"
)

var (
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
		return DeleteSessionOutput{}, fmt.Errorf("failed to extract user id from access token: %w", err)
	}

	userIDFromRefresh, err := jwt.ExtractID(input.RefreshToken)
	if err != nil {
		return DeleteSessionOutput{}, fmt.Errorf("failed to extract user id from refresh token: %w", err)
	}

	if userID != userIDFromRefresh {
		return DeleteSessionOutput{}, fmt.Errorf("invalid user id: %s != %s", userID, userIDFromRefresh)
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
