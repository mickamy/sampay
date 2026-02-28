package usecase_test

import (
	"testing"

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

func TestArchiveEvent_Do(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		ev := fixture.Event(func(e *model.Event) { e.UserID = endUser.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		sut := usecase.NewArchiveEvent(infra)
		out, err := sut.Do(ctx, usecase.ArchiveEventInput{ID: ev.ID})

		require.NoError(t, err)
		assert.NotNil(t, out.Event.ArchivedAt)

		// verify persisted
		persisted, err := query.Events(infra.ReaderDB).Where("id = ?", ev.ID).First(t.Context())
		require.NoError(t, err)
		assert.NotNil(t, persisted.ArchivedAt)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		sut := usecase.NewArchiveEvent(infra)
		_, err := sut.Do(ctx, usecase.ArchiveEventInput{ID: "nonexistent"})

		require.Error(t, err)
		require.ErrorIs(t, err, usecase.ErrArchiveEventNotFound)
	})

	t.Run("forbidden", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		owner := tseed.EndUser(t, infra.WriterDB)
		other := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), other.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		ev := fixture.Event(func(e *model.Event) { e.UserID = owner.UserID })
		require.NoError(t, query.Events(infra.WriterDB).Create(t.Context(), &ev))

		sut := usecase.NewArchiveEvent(infra)
		_, err := sut.Do(ctx, usecase.ArchiveEventInput{ID: ev.ID})

		require.Error(t, err)
		require.ErrorIs(t, err, usecase.ErrArchiveEventForbidden)
	})
}
