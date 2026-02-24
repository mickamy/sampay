package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/user/fixture"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/domain/user/usecase"
)

func TestCheckSlugAvailability_Do(t *testing.T) {
	t.Parallel()

	t.Run("available when slug not taken", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)

		// act
		sut := usecase.NewCheckSlugAvailability(infra)
		out, err := sut.Do(t.Context(), usecase.CheckSlugAvailabilityInput{Slug: "fresh-slug"})

		// assert
		require.NoError(t, err)
		assert.True(t, out.Available)
	})

	t.Run("unavailable when slug already taken", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))

		// act
		sut := usecase.NewCheckSlugAvailability(infra)
		out, err := sut.Do(t.Context(), usecase.CheckSlugAvailabilityInput{Slug: endUser.Slug})

		// assert
		require.NoError(t, err)
		assert.False(t, out.Available)
	})

	t.Run("unavailable when slug is invalid", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)

		// act
		sut := usecase.NewCheckSlugAvailability(infra)
		out, err := sut.Do(t.Context(), usecase.CheckSlugAvailabilityInput{Slug: "AB"})

		// assert
		require.NoError(t, err)
		assert.False(t, out.Available)
	})
}
