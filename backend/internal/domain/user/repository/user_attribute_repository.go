package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/user/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UserAttribute interface {
	Create(ctx context.Context, m *model.UserAttribute) error
	Find(ctx context.Context, id string, scopes ...database.Scope) (*model.UserAttribute, error)
	Update(ctx context.Context, m *model.UserAttribute) error
	WithTx(tx *database.DB) UserAttribute
}

type userAttribute struct {
	db *database.DB
}

func NewUserAttribute(db *database.DB) UserAttribute {
	return &userAttribute{db: db}
}

func (repo *userAttribute) Create(ctx context.Context, m *model.UserAttribute) error {
	return repo.db.WithContext(ctx).Create(&m).Error
}

func (repo *userAttribute) Find(ctx context.Context, id string, scopes ...database.Scope) (*model.UserAttribute, error) {
	var m model.UserAttribute
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).Where("user_id = ?", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &m, err
}

func (repo *userAttribute) Update(ctx context.Context, m *model.UserAttribute) error {
	return repo.db.WithContext(ctx).Clauses(clause.Returning{}).Where("user_id = ?", m.UserID).Save(m).Error
}

func (repo *userAttribute) WithTx(tx *database.DB) UserAttribute {
	return &userAttribute{db: tx}
}
