package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/event/fixture"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/query"
	"github.com/mickamy/sampay/internal/domain/event/usecase"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func TestJoinEvent_Do(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)

		ev := fixture.Event(func(e *model.Event) {
			e.UserID = endUser.UserID
			e.TierCount = 3
		})
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		t2 := fixture.EventTier(func(m *model.EventTier) {
			m.EventID = ev.ID
			m.Tier = 2
			m.Count = 3
			m.Amount = 5000
		})
		require.NoError(t, query.EventTiers(infra.WriterDB).CreateAll(t.Context(), []*model.EventTier{&t2}))

		sut := usecase.NewJoinEvent(infra)
		out, err := sut.Do(t.Context(), usecase.JoinEventInput{
			EventID: ev.ID,
			Name:    "Alice",
			Tier:    2,
		})

		require.NoError(t, err)
		assert.Equal(t, ev.ID, out.Participant.EventID)
		assert.Equal(t, "Alice", out.Participant.Name)
		assert.Equal(t, 2, out.Participant.Tier)
		assert.Equal(t, 5000, out.Participant.Amount)
		assert.Equal(t, model.ParticipantStatusUnpaid, out.Participant.Status)
	})

	t.Run("event not found", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)

		sut := usecase.NewJoinEvent(infra)
		_, err := sut.Do(t.Context(), usecase.JoinEventInput{
			EventID: "nonexistent",
			Name:    "Alice",
			Tier:    1,
		})

		require.ErrorIs(t, err, usecase.ErrJoinEventNotFound)
	})

	t.Run("empty name", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)

		sut := usecase.NewJoinEvent(infra)
		_, err := sut.Do(t.Context(), usecase.JoinEventInput{
			EventID: "any",
			Name:    "",
			Tier:    1,
		})

		require.ErrorIs(t, err, usecase.ErrJoinEventEmptyName)
	})

	t.Run("invalid tier", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)

		ev := fixture.Event(func(e *model.Event) {
			e.UserID = endUser.UserID
			e.TierCount = 3
		})
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		sut := usecase.NewJoinEvent(infra)
		_, err := sut.Do(t.Context(), usecase.JoinEventInput{
			EventID: ev.ID,
			Name:    "Bob",
			Tier:    4, // > tier_count(3)
		})

		require.ErrorIs(t, err, usecase.ErrJoinEventInvalidTier)
	})

	t.Run("archived event", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)

		now := time.Now()
		ev := fixture.Event(func(e *model.Event) {
			e.UserID = endUser.UserID
			e.ArchivedAt = &now
		})
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		sut := usecase.NewJoinEvent(infra)
		_, err := sut.Do(t.Context(), usecase.JoinEventInput{
			EventID: ev.ID,
			Name:    "Alice",
			Tier:    1,
		})

		require.ErrorIs(t, err, usecase.ErrJoinEventArchived)
	})
}
