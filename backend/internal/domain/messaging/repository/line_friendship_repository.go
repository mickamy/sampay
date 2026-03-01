package repository

import (
	"context"
	"fmt"

	"github.com/mickamy/sampay/internal/domain/messaging/model"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

type LineFriendship interface {
	Upsert(ctx context.Context, m *model.LineFriendship) error
	WithTx(tx *database.DB) LineFriendship
}

type lineFriendship struct {
	db *database.DB
}

func NewLineFriendship(db *database.DB) LineFriendship {
	return &lineFriendship{db: db}
}

func (repo *lineFriendship) Upsert(ctx context.Context, m *model.LineFriendship) error {
	_, err := repo.db.ExecContext(ctx,
		`INSERT INTO line_friendships (end_user_id, is_friend, updated_at)
		 VALUES ($1, $2, CURRENT_TIMESTAMP)
		 ON CONFLICT (end_user_id) DO UPDATE SET is_friend = $2, updated_at = CURRENT_TIMESTAMP`,
		m.EndUserID, m.IsFriend,
	)
	if err != nil {
		return fmt.Errorf("repository: failed to upsert line friendship: %w", err)
	}
	return nil
}

func (repo *lineFriendship) WithTx(tx *database.DB) LineFriendship {
	return &lineFriendship{db: tx}
}
