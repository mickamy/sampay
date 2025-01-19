package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/ptr"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type UserProfile interface {
	Create(ctx context.Context, m *model.UserProfile) error
	Find(ctx context.Context, id string, scopes ...database.Scope) (*model.UserProfile, error)
	FindBySlug(ctx context.Context, slug string, scopes ...database.Scope) (*model.UserProfile, error)
	Update(ctx context.Context, m *model.UserProfile) error
	WithTx(tx *database.DB) UserProfile
}

type userProfile struct {
	db *database.DB
}

func NewUserProfile(db *database.DB) UserProfile {
	return &userProfile{db: db}
}

func (repo *userProfile) Create(ctx context.Context, m *model.UserProfile) error {
	return repo.db.WithContext(ctx).Create(&m).Error
}

func (repo *userProfile) Find(ctx context.Context, id string, scopes ...database.Scope) (*model.UserProfile, error) {
	var m model.UserProfile
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		First(&m, "user_id = ?", id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (repo *userProfile) FindBySlug(ctx context.Context, slug string, scopes ...database.Scope) (*model.UserProfile, error) {
	var m model.UserProfile
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		Joins("LEFT OUTER JOIN users on user_profiles.user_id = users.id").
		First(&m, "users.slug = ?", slug).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (repo *userProfile) Update(ctx context.Context, m *model.UserProfile) error {
	return repo.db.WithContext(ctx).
		Where("user_id = ?", m.UserID).
		Updates(m).
		Update("bio", ptr.NullIfZero(m.Bio)).
		Update("image_id", ptr.NullIfZero(m.ImageID)).
		Debug().
		Error
}

func (repo *userProfile) WithTx(tx *database.DB) UserProfile {
	return &userProfile{db: tx}
}

func UserProfileJoinPreloadImage(db *database.DB) *database.DB {
	grm := db.Preload("Image")
	return &database.DB{DB: grm}
}
