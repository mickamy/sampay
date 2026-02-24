package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/user/fixture"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/domain/user/repository"
)

func TestUserPaymentMethod_CreateAll(t *testing.T) {
	t.Parallel()

	// arrange
	db := newReadWriter(t)
	user := fixture.User(nil)
	require.NoError(t, query.Users(db.Writer.DB).Create(t.Context(), &user))
	endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
	require.NoError(t, query.EndUsers(db.Writer.DB).Create(t.Context(), &endUser))

	methods := []model.UserPaymentMethod{
		fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) {
			m.UserID = user.ID
			m.Type = "paypay"
			m.DisplayOrder = 0
		}),
		fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) {
			m.UserID = user.ID
			m.Type = "kyash"
			m.DisplayOrder = 1
		}),
	}

	// act
	sut := repository.NewUserPaymentMethod(db.Writer.DB)
	err := sut.CreateAll(t.Context(), methods)

	// assert
	require.NoError(t, err)
	got, err := query.UserPaymentMethods(db.Reader.DB).Where("user_id = ?", user.ID).All(t.Context())
	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestUserPaymentMethod_ListByUserID(t *testing.T) {
	t.Parallel()

	t.Run("returns methods ordered by display_order", func(t *testing.T) {
		t.Parallel()

		// arrange
		db := newReadWriter(t)
		user := fixture.User(nil)
		require.NoError(t, query.Users(db.Writer.DB).Create(t.Context(), &user))
		endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
		require.NoError(t, query.EndUsers(db.Writer.DB).Create(t.Context(), &endUser))

		m1 := fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) {
			m.UserID = user.ID
			m.Type = "kyash"
			m.DisplayOrder = 1
		})
		m2 := fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) {
			m.UserID = user.ID
			m.Type = "paypay"
			m.DisplayOrder = 0
		})
		require.NoError(t, query.UserPaymentMethods(db.Writer.DB).Create(t.Context(), &m1))
		require.NoError(t, query.UserPaymentMethods(db.Writer.DB).Create(t.Context(), &m2))

		// act
		sut := repository.NewUserPaymentMethod(db.Reader.DB)
		got, err := sut.ListByUserID(t.Context(), user.ID)

		// assert
		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, "paypay", got[0].Type)
		assert.Equal(t, "kyash", got[1].Type)
	})

	t.Run("returns empty when no methods", func(t *testing.T) {
		t.Parallel()

		// arrange
		db := newReadWriter(t)

		// act
		sut := repository.NewUserPaymentMethod(db.Reader.DB)
		got, err := sut.ListByUserID(t.Context(), "nonexistent")

		// assert
		require.NoError(t, err)
		assert.Empty(t, got)
	})
}

func TestUserPaymentMethod_DeleteByUserID(t *testing.T) {
	t.Parallel()

	// arrange
	db := newReadWriter(t)
	user := fixture.User(nil)
	require.NoError(t, query.Users(db.Writer.DB).Create(t.Context(), &user))
	endUser := fixture.EndUser(func(m *model.EndUser) { m.UserID = user.ID })
	require.NoError(t, query.EndUsers(db.Writer.DB).Create(t.Context(), &endUser))

	m := fixture.UserPaymentMethod(func(m *model.UserPaymentMethod) { m.UserID = user.ID })
	require.NoError(t, query.UserPaymentMethods(db.Writer.DB).Create(t.Context(), &m))

	// act
	sut := repository.NewUserPaymentMethod(db.Writer.DB)
	err := sut.DeleteByUserID(t.Context(), user.ID)

	// assert
	require.NoError(t, err)
	got, err := query.UserPaymentMethods(db.Reader.DB).Where("user_id = ?", user.ID).All(t.Context())
	require.NoError(t, err)
	assert.Empty(t, got)
}
