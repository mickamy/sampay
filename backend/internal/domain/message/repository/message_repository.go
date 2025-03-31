package repository

import (
	"context"

	"mickamy.com/sampay/internal/domain/message/model"
	"mickamy.com/sampay/internal/infra/storage/database"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type Message interface {
	Create(ctx context.Context, m *model.Message) error
	ListByReceiverID(ctx context.Context, receiverID string, scopes ...database.Scope) ([]model.Message, error)
	WithTx(tx *database.DB) Message
}

type message struct {
	db *database.DB
}

func NewMessage(db *database.DB) Message {
	return &message{db: db}
}

func (repo *message) Create(ctx context.Context, m *model.Message) error {
	return repo.db.WithContext(ctx).Create(m).Error
}

func (repo *message) ListByReceiverID(ctx context.Context, receiverID string, scopes ...database.Scope) ([]model.Message, error) {
	var messages []model.Message
	err := repo.db.WithContext(ctx).
		Scopes(database.Scopes(scopes).Gorm()...).
		Where("receiver_id = ?", receiverID).
		Find(&messages).Error
	return messages, err
}

func (repo *message) WithTx(tx *database.DB) Message {
	return &message{db: tx}
}

func MessageCreatedAtDesc(db *database.DB) *database.DB {
	return &database.DB{DB: db.Order("created_at DESC")}
}
