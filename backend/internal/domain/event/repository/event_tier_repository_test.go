package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/event/fixture"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/repository"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func TestEventTierRepository(t *testing.T) {
	t.Parallel()

	t.Run("create and list", func(t *testing.T) {
		t.Parallel()

		db := newReadWriter(t)
		endUser := tseed.EndUser(t, db.Writer)
		ev := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
		require.NoError(t, repository.NewEvent(db.Writer.DB).Create(t.Context(), &ev))

		repo := repository.NewEventTier(db.Writer.DB)

		t1 := fixture.EventTier(func(m *model.EventTier) {
			m.EventID = ev.ID
			m.Tier = 3
			m.Count = 2
			m.Amount = 5000
		})
		t2 := fixture.EventTier(func(m *model.EventTier) {
			m.EventID = ev.ID
			m.Tier = 1
			m.Count = 3
			m.Amount = 2000
		})
		require.NoError(t, repo.CreateAll(t.Context(), []*model.EventTier{&t1, &t2}))

		tiers, err := repo.ListByEventID(t.Context(), ev.ID)
		require.NoError(t, err)
		require.Len(t, tiers, 2)
		// ordered by tier ASC
		assert.Equal(t, 1, tiers[0].Tier)
		assert.Equal(t, 3, tiers[1].Tier)
	})

	t.Run("delete by event id", func(t *testing.T) {
		t.Parallel()

		db := newReadWriter(t)
		endUser := tseed.EndUser(t, db.Writer)
		ev := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
		require.NoError(t, repository.NewEvent(db.Writer.DB).Create(t.Context(), &ev))

		repo := repository.NewEventTier(db.Writer.DB)

		tier := fixture.EventTier(func(m *model.EventTier) {
			m.EventID = ev.ID
			m.Tier = 1
			m.Count = 5
			m.Amount = 1000
		})
		require.NoError(t, repo.CreateAll(t.Context(), []*model.EventTier{&tier}))

		require.NoError(t, repo.DeleteByEventID(t.Context(), ev.ID))

		tiers, err := repo.ListByEventID(t.Context(), ev.ID)
		require.NoError(t, err)
		assert.Empty(t, tiers)
	})
}
