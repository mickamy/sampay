package repository

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
	"mickamy.com/sampay/internal/domain/auth/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type Session interface {
	Create(ctx context.Context, session model.Session) error

	Delete(ctx context.Context, session model.Session) error

	AccessTokenExists(ctx context.Context, userID string, accessToken string) (bool, error)

	RefreshTokenExists(ctx context.Context, userID string, refreshToken string) (bool, error)
}

type session struct {
	kvs *kvs.KVS
}

func NewSession(kvs *kvs.KVS) Session {
	return &session{kvs: kvs}
}

func (repo *session) Create(ctx context.Context, session model.Session) error {
	if session.UserID == "" || session.Tokens.Access.Value == "" || session.Tokens.Refresh.Value == "" {
		return fmt.Errorf("session contains empty value")
	}
	if err := repo.kvs.Set(
		ctx,
		repo.accessTokenKey(session.UserID, session.Tokens.Access.Value),
		session.Tokens.Access.Value,
		session.Tokens.Access.Expiration(),
	).Err(); err != nil {
		return fmt.Errorf("failed to set access token: %w", err)
	}
	if err := repo.kvs.Set(
		ctx,
		repo.refreshTokenKey(session.UserID, session.Tokens.Refresh.Value),
		session.Tokens.Refresh.Value,
		session.Tokens.Refresh.Expiration(),
	).Err(); err != nil {
		return fmt.Errorf("failed to set refresh token: %w", err)
	}

	return nil
}

func (repo *session) Delete(ctx context.Context, session model.Session) error {
	if err := repo.kvs.Del(ctx, repo.accessTokenKey(session.UserID, session.Tokens.Access.Value)).Err(); err != nil {
		return fmt.Errorf("failed to delete access token: %w", err)
	}
	if err := repo.kvs.Del(ctx, repo.refreshTokenKey(session.UserID, session.Tokens.Refresh.Value)).Err(); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}

func (repo *session) AccessTokenExists(ctx context.Context, userID string, accessToken string) (bool, error) {
	exists, err := repo.kvs.Exists(ctx, repo.accessTokenKey(userID, accessToken)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check access token existence: %w", err)
	}
	return exists == 1, nil
}

func (repo *session) RefreshTokenExists(ctx context.Context, userID string, refreshToken string) (bool, error) {
	exists, err := repo.kvs.Exists(ctx, repo.refreshTokenKey(userID, refreshToken)).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check refresh token existence: %w", err)
	}
	return exists == 1, nil
}

func (repo *session) accessTokenKey(userID, accessToken string) string {
	return fmt.Sprintf("session:%s:access_token:%s", userID, accessToken)
}

func (repo *session) refreshTokenKey(userID, refreshToken string) string {
	return fmt.Sprintf("session:%s:refresh_token:%s", userID, refreshToken)
}
