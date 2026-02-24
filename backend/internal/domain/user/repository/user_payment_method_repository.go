package repository

import (
	"context"
	"fmt"

	"github.com/mickamy/ormgen/scope"

	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

type UserPaymentMethod interface {
	CreateAll(ctx context.Context, methods []*model.UserPaymentMethod) error
	ListByUserID(ctx context.Context, userID string, scopes ...scope.Scope) ([]model.UserPaymentMethod, error)
	DeleteByUserID(ctx context.Context, userID string) error
	WithTx(tx *database.DB) UserPaymentMethod
}

type userPaymentMethod struct {
	db *database.DB
}

func NewUserPaymentMethod(db *database.DB) UserPaymentMethod {
	return &userPaymentMethod{db: db}
}

func (repo *userPaymentMethod) CreateAll(ctx context.Context, methods []*model.UserPaymentMethod) error {
	if err := query.UserPaymentMethods(repo.db).CreateAll(ctx, methods); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *userPaymentMethod) ListByUserID(
	ctx context.Context, userID string, scopes ...scope.Scope,
) ([]model.UserPaymentMethod, error) {
	methods, err := query.UserPaymentMethods(repo.db).
		Scopes(scopes...).
		Where("user_id = ?", userID).
		OrderBy("display_order ASC").
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}
	return methods, nil
}

func (repo *userPaymentMethod) DeleteByUserID(ctx context.Context, userID string) error {
	if err := query.UserPaymentMethods(repo.db).Where("user_id = ?", userID).Delete(ctx); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *userPaymentMethod) WithTx(tx *database.DB) UserPaymentMethod {
	return &userPaymentMethod{db: tx}
}
