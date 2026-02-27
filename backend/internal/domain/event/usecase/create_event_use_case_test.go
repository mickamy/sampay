package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/mickamy/sampay/internal/domain/event/usecase"
	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/test/tseed"
)

func TestCreateEvent_Do(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		sut := usecase.NewCreateEvent(infra)
		out, err := sut.Do(ctx, usecase.CreateEventInput{
			Title:       "party",
			Description: "new year party",
			TotalAmount: 30000,
			TierCount:   3,
			HeldAt:      time.Now().Add(24 * time.Hour),
		})

		require.NoError(t, err)
		assert.NotEmpty(t, out.Event.ID)
		assert.Equal(t, endUser.UserID, out.Event.UserID)
		assert.Equal(t, "party", out.Event.Title)
		assert.Equal(t, 30000, out.Event.TotalAmount)
		assert.Equal(t, 3, out.Event.TierCount)
	})

	t.Run("empty title", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		sut := usecase.NewCreateEvent(infra)
		_, err := sut.Do(ctx, usecase.CreateEventInput{
			Title:       "",
			TotalAmount: 30000,
			TierCount:   1,
			HeldAt:      time.Now().Add(24 * time.Hour),
		})

		require.Error(t, err)
		require.ErrorIs(t, err, usecase.ErrValidateEventEmptyTitle)
	})

	t.Run("negative total amount", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		sut := usecase.NewCreateEvent(infra)
		_, err := sut.Do(ctx, usecase.CreateEventInput{
			Title:       "party",
			TotalAmount: -30000,
			TierCount:   1,
			HeldAt:      time.Now().Add(24 * time.Hour),
		})

		require.Error(t, err)
		require.ErrorIs(t, err, usecase.ErrValidateEventNonPositiveTotalAmount)
	})

	t.Run("invalid tier_count", func(t *testing.T) {
		t.Parallel()

		infra := newInfra(t)
		endUser := tseed.EndUser(t, infra.WriterDB)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), endUser.UserID)
		ctx = contexts.SetLanguage(ctx, language.Japanese)

		sut := usecase.NewCreateEvent(infra)
		_, err := sut.Do(ctx, usecase.CreateEventInput{
			Title:       "party",
			TotalAmount: 30000,
			TierCount:   2,
			HeldAt:      time.Now().Add(24 * time.Hour),
		})

		require.Error(t, err)
		require.ErrorIs(t, err, usecase.ErrValidateEventInvalidTierCount)
	})
}
