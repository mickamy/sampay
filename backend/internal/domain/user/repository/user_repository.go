package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/user/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type User interface {
	Create(ctx context.Context, m *model.User) error
	FindByID(ctx context.Context, id string, scopes ...database.Scope) (*model.User, error)
	FindBySlug(ctx context.Context, slug string, scopes ...database.Scope) (*model.User, error)
	FindByEmail(ctx context.Context, email string, scopes ...database.Scope) (*model.User, error)
	FindByEmailOrSlug(ctx context.Context, emailOrSlug string, scopes ...database.Scope) (*model.User, error)
	Upsert(ctx context.Context, m *model.User) error
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

func (repo *user) FindByID(ctx context.Context, id string, scopes ...database.Scope) (*model.User, error) {
	var m model.User
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		First(&m, "id = ?", id).
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

func (repo *user) FindByEmail(ctx context.Context, emailOrSlug string, scopes ...database.Scope) (*model.User, error) {
	var m model.User
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		Joins("LEFT OUTER JOIN authentications ON users.id = authentications.user_id").
		Where("(authentications.identifier = ? AND authentications.type = ?)", emailOrSlug, authModel.AuthenticationTypeEmailPassword).
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

func (repo *user) FindByEmailOrSlug(ctx context.Context, emailOrSlug string, scopes ...database.Scope) (*model.User, error) {
	var m model.User
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).
		Joins("LEFT OUTER JOIN authentications ON users.id = authentications.user_id").
		Where("(authentications.identifier = ? AND authentications.type = ?) OR users.slug = ?", emailOrSlug, authModel.AuthenticationTypeEmailPassword, emailOrSlug).
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

func (repo *user) Upsert(ctx context.Context, m *model.User) error {
	if m.Slug == "" {
		return errors.New("slug is required")
	}

	var id string
	err := repo.db.WithContext(ctx).Model(&model.User{}).
		Where("slug = ?", m.Slug).
		Limit(1).
		Pluck("id", &id).
		Error
	if err != nil {
		return fmt.Errorf("failed to find user by slug: %w", err)
	}
	if id != "" {
		return repo.db.WithContext(ctx).Model(m).Clauses(clause.Returning{}).Where("slug = ?").Updates(m).Error
	}
	return repo.db.WithContext(ctx).Create(m).Error
}

func (repo *user) WithTx(tx *database.DB) User {
	return &user{db: tx}
}
