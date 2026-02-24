package usecase_test

import (
	"testing"

	"github.com/mickamy/errx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/user/fixture"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/domain/user/usecase"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

func TestUpdateSlug_Do(t *testing.T) {
	t.Parallel()

	t.Run("updates slug successfully", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))
		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewUpdateSlug(infra)
		out, err := sut.Do(ctx, usecase.UpdateSlugInput{Slug: "new-slug"})

		// assert
		require.NoError(t, err)
		assert.Equal(t, "new-slug", out.User.Slug)
	})

	t.Run("returns error when slug is invalid", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		ctx := contexts.SetAuthenticatedUserID(t.Context(), "dummy")

		// act
		sut := usecase.NewUpdateSlug(infra)
		_, err := sut.Do(ctx, usecase.UpdateSlugInput{Slug: "AB"})

		// assert
		require.Error(t, err)
		var ex *errx.Error
		require.ErrorAs(t, err, &ex)
		assert.Equal(t, errx.InvalidArgument, ex.Code())
	})

	t.Run("returns error when slug already taken", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user1 := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user1))
		endUser1 := fixture.EndUser(func(m *model.EndUser) {
			m.UserID = user1.ID
			m.Slug = "taken-slug"
		})
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser1))

		user2 := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user2))
		endUser2 := fixture.EndUser(func(m *model.EndUser) { m.UserID = user2.ID })
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser2))
		ctx := contexts.SetAuthenticatedUserID(t.Context(), user2.ID)

		// act
		sut := usecase.NewUpdateSlug(infra)
		_, err := sut.Do(ctx, usecase.UpdateSlugInput{Slug: "taken-slug"})

		// assert
		require.Error(t, err)
		var ex *errx.Error
		require.ErrorAs(t, err, &ex)
		assert.Equal(t, errx.AlreadyExists, ex.Code())
	})
}
