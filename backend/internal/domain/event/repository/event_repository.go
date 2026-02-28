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

type Event interface {
	Create(ctx context.Context, m *model.Event) error
	Get(ctx context.Context, id string, scopes ...scope.Scope) (model.Event, error)
	ListByUserID(ctx context.Context, userID string, scopes ...scope.Scope) ([]model.Event, error)
	Update(ctx context.Context, m *model.Event) error
	Delete(ctx context.Context, id string) error
	WithTx(tx *database.DB) Event
}

type event struct {
	db *database.DB
}

func NewEvent(db *database.DB) Event {
	return &event{db: db}
}

func (repo *event) Create(ctx context.Context, m *model.Event) error {
	if err := query.Events(repo.db).Create(ctx, m); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *event) Get(ctx context.Context, id string, scopes ...scope.Scope) (model.Event, error) {
	m, err := query.Events(repo.db).Scopes(scopes...).Where("id = ?", id).First(ctx)
	if errors.Is(err, orm.ErrNotFound) {
		return model.Event{}, database.ErrNotFound
	}
	if err != nil {
		return m, fmt.Errorf("repository: %w", err)
	}
	m.SortTiers()
	return m, nil
}

func (repo *event) ListByUserID(ctx context.Context, userID string, scopes ...scope.Scope) ([]model.Event, error) {
	events, err := query.Events(repo.db).
		Scopes(scopes...).
		Where("user_id = ?", userID).
		OrderBy("held_at DESC").
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}
	for i := range events {
		events[i].SortTiers()
	}
	return events, nil
}

func (repo *event) Update(ctx context.Context, m *model.Event) error {
	if err := query.Events(repo.db).Update(ctx, m); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *event) Delete(ctx context.Context, id string) error {
	if err := query.Events(repo.db).Where("id = ?", id).Delete(ctx); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *event) WithTx(tx *database.DB) Event {
	return &event{db: tx}
}

func EventPreloadParticipants() scope.Scope {
	return scope.Preload("Participants")
}

func EventPreloadTiers() scope.Scope {
	return scope.Preload("Tiers")
}
