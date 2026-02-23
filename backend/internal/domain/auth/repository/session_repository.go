package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/infra/storage/kvs"
	"github.com/mickamy/sampay/internal/lib/logger"
)

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
		return errors.New("session contains empty value")
	}
	if err := repo.kvs.Set(
		ctx,
		repo.accessTokenKey(session.UserID, session.Tokens.Access.Value),
		session.Tokens.Access.Value,
		session.Tokens.Access.Expiration(),
	); err != nil {
		return fmt.Errorf("repository: failed to set access token: %w", err)
	}
	if err := repo.kvs.Set(
		ctx,
		repo.refreshTokenKey(session.UserID, session.Tokens.Refresh.Value),
		session.Tokens.Refresh.Value,
		session.Tokens.Refresh.Expiration(),
	); err != nil {
		return fmt.Errorf("repository: failed to set refresh token: %w", err)
	}

	logger.Info(ctx, "session created", "user", session.UserID, "session", session)

	return nil
}

func (repo *session) Delete(ctx context.Context, session model.Session) error {
	if err := repo.kvs.Del(ctx, repo.accessTokenKey(session.UserID, session.Tokens.Access.Value)); err != nil {
		return fmt.Errorf("repository: failed to delete access token: %w", err)
	}
	if err := repo.kvs.Del(ctx, repo.refreshTokenKey(session.UserID, session.Tokens.Refresh.Value)); err != nil {
		return fmt.Errorf("repository: failed to delete refresh token: %w", err)
	}

	return nil
}

func (repo *session) AccessTokenExists(ctx context.Context, userID string, accessToken string) (bool, error) {
	exists, err := repo.kvs.Exists(ctx, repo.accessTokenKey(userID, accessToken))
	if err != nil {
		return false, fmt.Errorf("failed to check access token existence: %w", err)
	}
	return exists, nil
}

func (repo *session) RefreshTokenExists(ctx context.Context, userID string, refreshToken string) (bool, error) {
	exists, err := repo.kvs.Exists(ctx, repo.refreshTokenKey(userID, refreshToken))
	if err != nil {
		return false, fmt.Errorf("failed to check refresh token existence: %w", err)
	}
	return exists, nil
}

func (repo *session) accessTokenKey(userID, accessToken string) string {
	return fmt.Sprintf("session:%s:access:%s", userID, accessToken)
}

func (repo *session) refreshTokenKey(userID, refreshToken string) string {
	return fmt.Sprintf("session:%s:refresh:%s", userID, refreshToken)
}
