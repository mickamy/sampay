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
)

func TestGetUserProfile_Do(t *testing.T) {
	t.Parallel()

	t.Run("returns user and payment methods", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))
		pm := fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) { m.UserID = user.ID })
		require.NoError(t, query.UserPaymentMethods(infra.WriterDB).Create(t.Context(), &pm))

		// act
		sut := usecase.NewGetUserProfile(infra)
		out, err := sut.Do(t.Context(), usecase.GetUserProfileInput{Slug: endUser.Slug})

		// assert
		require.NoError(t, err)
		assert.Equal(t, endUser.UserID, out.User.UserID)
		assert.Equal(t, endUser.Slug, out.User.Slug)
		assert.Len(t, out.PaymentMethods, 1)
		assert.Equal(t, pm.ID, out.PaymentMethods[0].ID)
	})

	t.Run("returns not found for unknown slug", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)

		// act
		sut := usecase.NewGetUserProfile(infra)
		_, err := sut.Do(t.Context(), usecase.GetUserProfileInput{Slug: "nonexistent"})

		// assert
		require.Error(t, err)
		var ex *errx.Error
		require.ErrorAs(t, err, &ex)
		assert.Equal(t, errx.NotFound, ex.Code())
	})

	t.Run("returns empty payment methods when none registered", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))

		// act
		sut := usecase.NewGetUserProfile(infra)
		out, err := sut.Do(t.Context(), usecase.GetUserProfileInput{Slug: endUser.Slug})

		// assert
		require.NoError(t, err)
		assert.Equal(t, endUser.UserID, out.User.UserID)
		assert.Empty(t, out.PaymentMethods)
	})
}
