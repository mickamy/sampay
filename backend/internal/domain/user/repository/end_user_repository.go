package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mickamy/ormgen/orm"
	"github.com/mickamy/ormgen/scope"

	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

type EndUser interface {
	Create(ctx context.Context, m *model.EndUser) error
	Get(ctx context.Context, id string, scopes ...scope.Scope) (model.EndUser, error)
	WithTx(tx *database.DB) EndUser
}

type endUser struct {
	db *database.DB
}

func NewEndUser(db *database.DB) EndUser {
	return &endUser{db: db}
}

func (repo *endUser) Create(ctx context.Context, m *model.EndUser) error {
	if err := query.EndUsers(repo.db).Create(ctx, m); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *endUser) Get(ctx context.Context, id string, scopes ...scope.Scope) (model.EndUser, error) {
	m, err := query.EndUsers(repo.db).Scopes(scopes...).Where("user_id = ?", id).First(ctx)
	if errors.Is(err, orm.ErrNotFound) {
		return model.EndUser{}, database.ErrNotFound
	}
	return m, err
}

func (repo *endUser) WithTx(tx *database.DB) EndUser {
	return &endUser{db: tx}
}
