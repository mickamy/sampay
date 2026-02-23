package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/user/fixture"
	"github.com/mickamy/sampay/internal/domain/user/query"
	"github.com/mickamy/sampay/internal/domain/user/repository"
)

func TestUser_Create(t *testing.T) {
	t.Parallel()

	// arrange
	db := newReadWriter(t)
	m := fixture.User(nil)

	// act
	sut := repository.NewUser(db.Writer.DB)
	err := sut.Create(t.Context(), &m)

	// assert
	require.NoError(t, err)
	got, err := query.Users(db.Reader.DB).Where("id = ?", m.ID).First(t.Context())
	require.NoError(t, err)
	assert.Equal(t, m.ID, got.ID)
}
