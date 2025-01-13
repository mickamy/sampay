package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/user/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UserLink interface {
	Create(ctx context.Context, m *model.UserLink) error
	ListByUserID(ctx context.Context, userID string, scopes ...database.Scope) ([]model.UserLink, error)
	Find(ctx context.Context, id string, scopes ...database.Scope) (*model.UserLink, error)
	GetLastDisplayOrderByUserID(ctx context.Context, userID string) (int, error)
	Update(ctx context.Context, m *model.UserLink) error
	Delete(ctx context.Context, id string) error
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

func (repo *userLink) ListByUserID(ctx context.Context, userID string, scopes ...database.Scope) ([]model.UserLink, error) {
	var ms []model.UserLink
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).Find(&ms, "user_id = ?", userID).Error
	return ms, err
}

func (repo *userLink) Find(ctx context.Context, id string, scopes ...database.Scope) (*model.UserLink, error) {
	m := new(model.UserLink)
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).First(m, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m, err
}

func (repo *userLink) GetLastDisplayOrderByUserID(ctx context.Context, userID string) (int, error) {
	var m model.UserLink
	err := repo.db.WithContext(ctx).
		Joins("DisplayAttribute").
		Order(`"DisplayAttribute".display_order DESC`).
		First(&m, "user_id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return -1, nil
		}
		return -1, err
	}
	return m.DisplayAttribute.DisplayOrder, err
}

func (repo *userLink) Update(ctx context.Context, m *model.UserLink) error {
	return repo.db.WithContext(ctx).Save(m).Error
}

func (repo *userLink) Delete(ctx context.Context, id string) error {
	return repo.db.WithContext(ctx).Delete(&model.UserLink{}, "id = ?", id).Error
}

func (repo *userLink) WithTx(tx *database.DB) UserLink {
	return &userLink{db: tx}
}

func UserLinkJoinDisplayAttribute(db *database.DB) *database.DB {
	grm := db.Joins("DisplayAttribute").Order(`"DisplayAttribute".display_order`)
	return &database.DB{DB: grm}
}

func UserLinkPreloadQRCode(db *database.DB) *database.DB {
	grm := db.Preload("QRCode")
	return &database.DB{DB: grm}
}
