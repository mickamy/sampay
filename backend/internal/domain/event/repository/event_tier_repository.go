package repository

import (
	"context"
	"fmt"

	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/query"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

type EventTier interface {
	CreateAll(ctx context.Context, tiers []*model.EventTier) error
	ListByEventID(ctx context.Context, eventID string) ([]model.EventTier, error)
	DeleteByEventID(ctx context.Context, eventID string) error
	WithTx(tx *database.DB) EventTier
}

type eventTier struct {
	db *database.DB
}

func NewEventTier(db *database.DB) EventTier {
	return &eventTier{db: db}
}

func (repo *eventTier) CreateAll(ctx context.Context, tiers []*model.EventTier) error {
	if err := query.EventTiers(repo.db).CreateAll(ctx, tiers); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *eventTier) ListByEventID(
	ctx context.Context, eventID string,
) ([]model.EventTier, error) {
	tiers, err := query.EventTiers(repo.db).
		Where("event_id = ?", eventID).
		OrderBy("tier ASC").
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository: %w", err)
	}
	return tiers, nil
}

func (repo *eventTier) DeleteByEventID(ctx context.Context, eventID string) error {
	if err := query.EventTiers(repo.db).Where("event_id = ?", eventID).Delete(ctx); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *eventTier) WithTx(tx *database.DB) EventTier {
	return &eventTier{db: tx}
}
