package tseed

import (
	"testing"

	"github.com/stretchr/testify/require"

	ufixture "github.com/mickamy/sampay/internal/domain/user/fixture"
	umodel "github.com/mickamy/sampay/internal/domain/user/model"
	uquery "github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

func EndUser(t *testing.T, db *database.Writer) umodel.EndUser {
	t.Helper()
	user := ufixture.User(nil)
	m := ufixture.EndUser(func(m *umodel.EndUser) {
		m.UserID = user.ID
	})
	require.NoError(t, uquery.Users(db).Create(t.Context(), &user))
	require.NoError(t, uquery.EndUsers(db).Create(t.Context(), &m))
	return m
}
