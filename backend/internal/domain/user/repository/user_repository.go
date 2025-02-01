package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/user/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type User interface {
	Create(ctx context.Context, m *model.User) error
	Get(ctx context.Context, id string, scopes ...database.Scope) (model.User, error)
	Find(ctx context.Context, id string, scopes ...database.Scope) (*model.User, error)
	FindBySlug(ctx context.Context, slug string, scopes ...database.Scope) (*model.User, error)
	FindByEmail(ctx context.Context, email string, scopes ...database.Scope) (*model.User, error)
	FindByEmailOrSlug(ctx context.Context, emailOrSlug string, scopes ...database.Scope) (*model.User, error)
	Update(ctx context.Context, m *model.User) error
	WithTx(tx *database.DB) User
}

type user struct {
	db *database.DB
}

func NewUser(db *database.DB) User {
	return &user{db: db}
}

func (repo *user) Create(ctx context.Context, m *model.User) error {
	return repo.db.WithContext(ctx).Create(&m).Error
}

func (repo *user) Get(ctx context.Context, id string, scopes ...database.Scope) (model.User, error) {
	m, err := repo.Find(ctx, id, scopes...)
	if err != nil {
		return model.User{}, err
	}
	if m == nil {
		return model.User{}, database.ErrRecordNotFound
	}
	return *m, err
}

func (repo *user) Find(ctx context.Context, id string, scopes ...database.Scope) (*model.User, error) {
	var m model.User
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		First(&m, "users.id = ?", id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (repo *user) FindBySlug(ctx context.Context, slug string, scopes ...database.Scope) (*model.User, error) {
	var m model.User
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		First(&m, "slug = ?", slug).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (repo *user) FindByEmail(ctx context.Context, email string, scopes ...database.Scope) (*model.User, error) {
	var m model.User
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).First(&m, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (repo *user) FindByEmailOrSlug(ctx context.Context, emailOrSlug string, scopes ...database.Scope) (*model.User, error) {
	var m model.User
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		Where("email = ? OR users.slug = ?", emailOrSlug, emailOrSlug).
		First(&m).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, err
}

func (repo *user) Update(ctx context.Context, m *model.User) error {
	return repo.db.WithContext(ctx).Save(m).Error
}

func (repo *user) WithTx(tx *database.DB) User {
	return &user{db: tx}
}

func UserJoinAttribute(tx *database.DB) *database.DB {
	return &database.DB{DB: tx.Joins("Attribute")}
}

func UserJoinProfile(tx *database.DB) *database.DB {
	return &database.DB{DB: tx.Joins("Profile")}
}

func UserJoinProfileAndImage(tx *database.DB) *database.DB {
	return &database.DB{DB: tx.Joins("Profile.Image")}
}

func UserPreloadLinksQRCodeAndDisplayAttributes(tx *database.DB) *database.DB {
	return &database.DB{DB: tx.Preload("Links.DisplayAttribute").Preload("Links.QRCode")}
}
