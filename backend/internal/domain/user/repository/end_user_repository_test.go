package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/user/fixture"
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/domain/user/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

func TestEndUser_Create(t *testing.T) {
	t.Parallel()

	// arrange
	db := newReadWriter(t)
	user := fixture.User(nil)
	require.NoError(t, query.Users(db.Writer.DB).Create(t.Context(), &user))
	m := fixture.EndUser(func(m *model.EndUser) {
		m.UserID = user.ID
	})

	// act
	sut := repository.NewEndUser(db.Writer.DB)
	err := sut.Create(t.Context(), &m)

	// assert
	require.NoError(t, err)
	got, err := query.EndUsers(db.Reader.DB).Where("user_id = ?", m.UserID).First(t.Context())
	require.NoError(t, err)
	assert.Equal(t, m.UserID, got.UserID)
	assert.Equal(t, m.Slug, got.Slug)
}

func TestEndUser_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(t *testing.T, db *database.ReadWriter) model.EndUser
		assert  func(t *testing.T, got model.EndUser, err error)
	}{
		{
			name: "found",
			arrange: func(t *testing.T, db *database.ReadWriter) model.EndUser {
				user := fixture.User(nil)
				m := fixture.EndUser(func(m *model.EndUser) {
					m.UserID = user.ID
				})
				require.NoError(t, query.Users(db.Writer.DB).Create(t.Context(), &user))
				require.NoError(t, query.EndUsers(db.Writer.DB).Create(t.Context(), &m))
				return m
			},
			assert: func(t *testing.T, got model.EndUser, err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.UserID)
				assert.NotEmpty(t, got.Slug)
			},
		},
		{
			name: "not found",
			arrange: func(t *testing.T, db *database.ReadWriter) model.EndUser {
				return fixture.EndUser(nil)
			},
			assert: func(t *testing.T, got model.EndUser, err error) {
				require.ErrorIs(t, err, database.ErrNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			db := newReadWriter(t)
			m := tt.arrange(t, db)

			// act
			sut := repository.NewEndUser(db.Reader.DB)
			got, err := sut.Get(t.Context(), m.UserID)

			// assert
			tt.assert(t, got, err)
		})
	}
}
