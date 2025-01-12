package repository

import (
	"context"

	"gorm.io/gorm/clause"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/user/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UserLinkProvider interface {
	Upsert(ctx context.Context, m *model.UserLinkProvider) error
	WithTx(tx *database.DB) UserLinkProvider
}

type userLinkProvider struct {
	db *database.DB
}

func NewUserLinkProvider(db *database.DB) UserLinkProvider {
	return &userLinkProvider{db: db}
}

func (repo *userLinkProvider) Upsert(ctx context.Context, m *model.UserLinkProvider) error {
	return repo.db.WithContext(ctx).Clauses(clause.Returning{}).Save(m).Error
}

func (repo *userLinkProvider) WithTx(tx *database.DB) UserLinkProvider {
	return &userLinkProvider{db: tx}
}
