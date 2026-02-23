package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/auth/fixture"
	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/domain/auth/query"
	"github.com/mickamy/sampay/internal/domain/auth/repository"
	ufixture "github.com/mickamy/sampay/internal/domain/user/fixture"
	umodel "github.com/mickamy/sampay/internal/domain/user/model"
	uquery "github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

func TestOAuthAccount_Create(t *testing.T) {
	t.Parallel()

	// arrange
	db := newReadWriter(t)
	user := ufixture.User(nil)
	endUser := ufixture.EndUser(func(m *umodel.EndUser) {
		m.UserID = user.ID
	})
	m := fixture.OAuthAccount(func(m *model.OAuthAccount) {
		m.EndUserID = endUser.UserID
	})
	require.NoError(t, uquery.Users(db.Writer.DB).Create(t.Context(), &user))
	require.NoError(t, uquery.EndUsers(db.Writer.DB).Create(t.Context(), &endUser))

	// act
	sut := repository.NewOAuthAccount(db.Writer.DB)
	err := sut.Create(t.Context(), &m)

	// assert
	require.NoError(t, err)
	got, err := query.OAuthAccounts(db.Reader.DB).Where("end_user_id = ?", m.EndUserID).First(t.Context())
	require.NoError(t, err)
	assert.Equal(t, m.EndUserID, got.EndUserID)
	assert.Equal(t, m.Provider, got.Provider)
	assert.Equal(t, m.UID, got.UID)
}

func TestOAuthAccount_GetByProviderAndUID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		arrange func(t *testing.T, db *database.ReadWriter) model.OAuthAccount
		assert  func(t *testing.T, got model.OAuthAccount, err error)
	}{
		{
			name: "found",
			arrange: func(t *testing.T, db *database.ReadWriter) model.OAuthAccount {
				user := ufixture.User(nil)
				endUser := ufixture.EndUser(func(m *umodel.EndUser) {
					m.UserID = user.ID
				})
				m := fixture.OAuthAccount(func(m *model.OAuthAccount) {
					m.EndUserID = endUser.UserID
				})
				require.NoError(t, uquery.Users(db.Writer.DB).Create(t.Context(), &user))
				require.NoError(t, uquery.EndUsers(db.Writer.DB).Create(t.Context(), &endUser))
				require.NoError(t, query.OAuthAccounts(db.Writer.DB).Create(t.Context(), &m))
				return m
			},
			assert: func(t *testing.T, got model.OAuthAccount, err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.ID)
			},
		},
		{
			name: "not found",
			arrange: func(t *testing.T, db *database.ReadWriter) model.OAuthAccount {
				return fixture.OAuthAccount(nil)
			},
			assert: func(t *testing.T, got model.OAuthAccount, err error) {
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
			sut := repository.NewOAuthAccount(db.Reader.DB)
			got, err := sut.GetByProviderAndUID(t.Context(), model.OAuthProvider(m.Provider), m.UID)

			// assert
			tt.assert(t, got, err)
		})
	}
}
