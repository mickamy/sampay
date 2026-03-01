package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mickamy/ormgen/orm"
	"github.com/mickamy/ormgen/scope"

	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/domain/auth/query"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

type OAuthAccount interface {
	Create(ctx context.Context, m *model.OAuthAccount) error
	GetByProviderAndUID(
		ctx context.Context, provider model.OAuthProvider, uid string, scopes ...scope.Scope,
	) (model.OAuthAccount, error)
	GetUIDByEndUserIDAndProvider(
		ctx context.Context, endUserID string, provider model.OAuthProvider,
	) (string, error)
	WithTx(tx *database.DB) OAuthAccount
}

type oauthAccount struct {
	db *database.DB
}

func NewOAuthAccount(db *database.DB) OAuthAccount {
	return &oauthAccount{db: db}
}

func (repo *oauthAccount) Create(ctx context.Context, m *model.OAuthAccount) error {
	if err := query.OAuthAccounts(repo.db).Create(ctx, m); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *oauthAccount) GetByProviderAndUID(
	ctx context.Context, provider model.OAuthProvider, uid string, scopes ...scope.Scope,
) (model.OAuthAccount, error) {
	m, err := query.OAuthAccounts(repo.db).
		Scopes(scopes...).
		Where("provider = ? AND uid = ?", provider, uid).
		First(ctx)
	if errors.Is(err, orm.ErrNotFound) {
		return model.OAuthAccount{}, database.ErrNotFound
	}
	if err != nil {
		return m, fmt.Errorf("failed to find oauth account: %w", err)
	}
	return m, nil
}

func (repo *oauthAccount) GetUIDByEndUserIDAndProvider(
	ctx context.Context, endUserID string, provider model.OAuthProvider,
) (string, error) {
	m, err := query.OAuthAccounts(repo.db).
		Where("end_user_id = ? AND provider = ?", endUserID, provider).
		First(ctx)
	if errors.Is(err, orm.ErrNotFound) {
		return "", database.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("failed to find oauth account: %w", err)
	}
	return m.UID, nil
}

func (repo *oauthAccount) WithTx(tx *database.DB) OAuthAccount {
	return &oauthAccount{db: tx}
}
