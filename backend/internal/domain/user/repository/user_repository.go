package repository

import (
	"context"
	"fmt"

	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

type User interface {
	Create(ctx context.Context, m *model.User) error
	WithTx(tx *database.DB) User
}

type user struct {
	db *database.DB
}

func NewUser(db *database.DB) User {
	return &user{db: db}
}

func (repo *user) Create(ctx context.Context, m *model.User) error {
	if err := query.Users(repo.db).Create(ctx, m); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *user) WithTx(tx *database.DB) User {
	return &user{db: tx}
}
