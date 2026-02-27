package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mickamy/ormgen/orm"
	"github.com/mickamy/ormgen/scope"

	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/query"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

type EventParticipant interface {
	Create(ctx context.Context, m *model.EventParticipant) error
	Get(ctx context.Context, id string, scopes ...scope.Scope) (model.EventParticipant, error)
	ListByEventID(ctx context.Context, eventID string, scopes ...scope.Scope) ([]model.EventParticipant, error)
	Update(ctx context.Context, m *model.EventParticipant) error
	WithTx(tx *database.DB) EventParticipant
}

type eventParticipant struct {
	db *database.DB
}

func NewEventParticipant(db *database.DB) EventParticipant {
	return &eventParticipant{db: db}
}

func (repo *eventParticipant) Create(ctx context.Context, m *model.EventParticipant) error {
	if err := query.EventParticipants(repo.db).Create(ctx, m); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *eventParticipant) Get(ctx context.Context, id string, scopes ...scope.Scope) (model.EventParticipant, error) {
	m, err := query.EventParticipants(repo.db).Scopes(scopes...).Where("id = ?", id).First(ctx)
	if errors.Is(err, orm.ErrNotFound) {
		return model.EventParticipant{}, database.ErrNotFound
	}
	if err != nil {
		return m, fmt.Errorf("repository: %w", err)
	}
	return m, nil
}

func (repo *eventParticipant) ListByEventID(ctx context.Context, eventID string, scopes ...scope.Scope) ([]model.EventParticipant, error) {
	participants, err := query.EventParticipants(repo.db).
		Scopes(scopes...).
		Where("event_id = ?", eventID).
		OrderBy("created_at ASC").
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}
	return participants, nil
}

func (repo *eventParticipant) Update(ctx context.Context, m *model.EventParticipant) error {
	if err := query.EventParticipants(repo.db).Update(ctx, m); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *eventParticipant) WithTx(tx *database.DB) EventParticipant {
	return &eventParticipant{db: tx}
}
