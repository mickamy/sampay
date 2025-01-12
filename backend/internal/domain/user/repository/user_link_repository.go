package repository

import (
	"context"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/user/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UserLink interface {
	Create(ctx context.Context, m *model.UserLink) error
	WithTx(tx *database.DB) UserLink
}

type userLink struct {
	db *database.DB
}

func NewUserLink(db *database.DB) UserLink {
	return &userLink{db: db}
}

func (repo *userLink) Create(ctx context.Context, m *model.UserLink) error {
	return repo.db.WithContext(ctx).Create(m).Error
}

func (repo *userLink) WithTx(tx *database.DB) UserLink {
	return &userLink{db: tx}
}
