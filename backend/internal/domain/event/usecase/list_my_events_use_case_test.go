package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/mickamy/sampay/internal/domain/event/fixture"
	"github.com/mickamy/sampay/internal/domain/event/model"
	"github.com/mickamy/sampay/internal/domain/event/query"
	"github.com/mickamy/sampay/internal/domain/event/usecase"
	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func TestListMyEvents_Do(t *testing.T) {
	t.Parallel()

	t.Run("returns events for authenticated user", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		for range 3 {
			ev := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
			require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))
		}

		sut := usecase.NewListMyEvents(infra)
		out, err := sut.Do(ctx, usecase.ListMyEventsInput{})

		require.NoError(t, err)
		assert.Len(t, out.Events, 3)
	})

	t.Run("returns empty when no events", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		sut := usecase.NewListMyEvents(infra)
		out, err := sut.Do(ctx, usecase.ListMyEventsInput{})

		require.NoError(t, err)
		assert.Empty(t, out.Events)
	})

	t.Run("returns only active events by default", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		// active event
		active := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &active))

		// archived event
		now := time.Now()
		archived := fixture.Event(func(e *model.Event) {
			e.UserID = endUser.UserID
			e.ArchivedAt = &now
		})
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &archived))

		sut := usecase.NewListMyEvents(infra)
		out, err := sut.Do(ctx, usecase.ListMyEventsInput{IncludeArchived: false})

		require.NoError(t, err)
		assert.Len(t, out.Events, 1)
		assert.Equal(t, active.ID, out.Events[0].ID)
	})

	t.Run("returns only archived events when include_archived is true", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		// active event
		active := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &active))

		// archived event
		now := time.Now()
		archived := fixture.Event(func(e *model.Event) {
			e.UserID = endUser.UserID
			e.ArchivedAt = &now
		})
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &archived))

		sut := usecase.NewListMyEvents(infra)
		out, err := sut.Do(ctx, usecase.ListMyEventsInput{IncludeArchived: true})

		require.NoError(t, err)
		assert.Len(t, out.Events, 1)
		assert.Equal(t, archived.ID, out.Events[0].ID)
	})
}
