package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/user/fixture"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/domain/user/usecase"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

func TestGetMe_Do(t *testing.T) {
	t.Parallel()

	t.Run("returns current user", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) {
			m.UserID = user.ID
			m.Slug = "my-slug"
		})
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))
		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewGetMe(infra)
		out, err := sut.Do(ctx, usecase.GetMeInput{})

		// assert
		require.NoError(t, err)
		assert.Equal(t, user.ID, out.User.UserID)
		assert.Equal(t, "my-slug", out.User.Slug)
	})
}
