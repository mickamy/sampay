package repository

import (
	"context"

	"gorm.io/gorm/clause"

	"mickamy.com/sampay/internal/domain/registration/model"
	"mickamy.com/sampay/internal/infra/storage/database"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UsageCategory interface {
	List(ctx context.Context, scopes ...database.Scope) ([]model.UsageCategory, error)
	Upsert(ctx context.Context, m *model.UsageCategory) error
	WithTx(tx *database.DB) UsageCategory
}

type usageCategory struct {
	db *database.DB
}

func NewUsageCategory(db *database.DB) UsageCategory {
	return &usageCategory{db: db}
}

func (repo *usageCategory) List(ctx context.Context, scopes ...database.Scope) ([]model.UsageCategory, error) {
	var ms []model.UsageCategory
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		Find(&ms).
		Error
	return ms, err
}

func (repo *usageCategory) Upsert(ctx context.Context, m *model.UsageCategory) error {
	return repo.db.WithContext(ctx).Clauses(clause.Returning{}).Save(m).Error
}

func (repo *usageCategory) WithTx(tx *database.DB) UsageCategory {
	return &usageCategory{db: tx}
}
