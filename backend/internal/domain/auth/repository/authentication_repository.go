package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/auth/model"
)

type AuthenticationKey struct {
	UserID     string
	Type       model.AuthenticationType
	Identifier string
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type Authentication interface {
	Create(ctx context.Context, m *model.Authentication) error
	FindByTypeAndIdentifier(ctx context.Context, authType model.AuthenticationType, identifier string, scopes database.Scopes) (*model.Authentication, error)
	ExistsByTypeAndIdentifier(ctx context.Context, authType model.AuthenticationType, identifier string) (bool, error)
	Update(ctx context.Context, m *model.Authentication) error
	Upsert(ctx context.Context, m *model.Authentication, key AuthenticationKey) error
	WithTx(tx *database.DB) Authentication
}

type authentication struct {
	db *database.DB
}

func NewAuthentication(db *database.DB) Authentication {
	return &authentication{db: db}
}

func (repo *authentication) Create(ctx context.Context, m *model.Authentication) error {
	return repo.db.WithContext(ctx).Create(&m).Error
}

func (repo *authentication) FindByTypeAndIdentifier(ctx context.Context, authType model.AuthenticationType, identifier string, scopes database.Scopes) (*model.Authentication, error) {
	m := new(model.Authentication)
	err := repo.db.WithContext(ctx).Scopes(scopes.Gorm()...).
		First(m, "type = ? AND identifier = ?", authType, identifier).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m, err
}

func (repo *authentication) ExistsByTypeAndIdentifier(ctx context.Context, authType model.AuthenticationType, identifier string) (bool, error) {
	var id string
	err := repo.db.WithContext(ctx).
		Model(model.Authentication{}).
		Where("type = ? AND identifier = ?", authType, identifier).
		Limit(1).
		Pluck("id", &id).Error
	return id != "", err
}

func (repo *authentication) Update(ctx context.Context, m *model.Authentication) error {
	return repo.db.WithContext(ctx).Save(m).Error
}

func (repo *authentication) Upsert(ctx context.Context, m *model.Authentication, key AuthenticationKey) error {
	var id string
	if err := repo.db.WithContext(ctx).Model(&model.Authentication{}).
		Where("user_id = ? AND type = ? AND identifier = ?", key.UserID, key.Type, key.Identifier).
		Pluck("id", &id).
		Error; err != nil {
		return fmt.Errorf("failed to check authentication existence: %w", err)
	}
	if id != "" {
		return repo.db.WithContext(ctx).
			Clauses(&clause.Returning{}).
			Where("user_id = ? AND type = ? AND identifier = ?", key.UserID, key.Type, key.Identifier).
			Updates(&m).Error
	}

	return repo.db.WithContext(ctx).Create(&m).Error
}

func (repo *authentication) WithTx(tx *database.DB) Authentication {
	return &authentication{db: tx}
}
