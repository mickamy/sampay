package repository

import (
	"context"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/notification/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type Notification interface {
	Create(ctx context.Context, m *model.Notification) error
	ListByUserID(ctx context.Context, userID string, scopes ...database.Scope) ([]model.Notification, error)
	WithTx(tx *database.DB) Notification
}

type notification struct {
	db *database.DB
}

func NewNotification(db *database.DB) Notification {
	return &notification{db: db}
}

func (repo *notification) Create(ctx context.Context, m *model.Notification) error {
	return repo.db.WithContext(ctx).Create(m).Error
}

func (repo *notification) ListByUserID(ctx context.Context, userID string, scopes ...database.Scope) ([]model.Notification, error) {
	var notifications []model.Notification
	err := repo.db.WithContext(ctx).Scopes(database.Scopes(scopes).Gorm()...).Find(&notifications, "user_id = ?", userID).Error
	return notifications, err
}

func (repo *notification) WithTx(tx *database.DB) Notification {
	return &notification{db: tx}
}

func NotificationJoinReadStatus(db *database.DB) *database.DB {
	return &database.DB{DB: db.Joins("ReadStatus")}
}
