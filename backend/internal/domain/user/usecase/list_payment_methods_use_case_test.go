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

func TestListPaymentMethods_Do(t *testing.T) {
	t.Parallel()

	t.Run("returns payment methods for authenticated user", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
		require.NoError(t, query.EndUsers(infra.WriterDB).Create(t.Context(), &endUser))

		m := fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) { m.UserID = user.ID })
		require.NoError(t, query.UserPaymentMethods(infra.WriterDB).Create(t.Context(), &m))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewListPaymentMethods(infra)
		out, err := sut.Do(ctx, usecase.ListPaymentMethodsInput{})

		// assert
		require.NoError(t, err)
		assert.Len(t, out.PaymentMethods, 1)
		assert.Equal(t, m.ID, out.PaymentMethods[0].ID)
	})

	t.Run("returns empty when no methods", func(t *testing.T) {
		t.Parallel()

		// arrange
		infra := newInfra(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(infra.WriterDB).Create(t.Context(), &user))

		ctx := contexts.SetAuthenticatedUserID(t.Context(), user.ID)

		// act
		sut := usecase.NewListPaymentMethods(newInfra(t))
		out, err := sut.Do(ctx, usecase.ListPaymentMethodsInput{})

		// assert
		require.NoError(t, err)
		assert.Empty(t, out.PaymentMethods)
	})
}
