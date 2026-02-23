package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/domain/auth/repository"
	"github.com/mickamy/sampay/internal/lib/ulid"
	"github.com/mickamy/sampay/internal/test/itest"
)

func TestSession_Create(t *testing.T) {
	t.Parallel()

	// arrange
	kvStore := itest.NewKVS(t)
	session := model.MustNewSession(ulid.New())
	sut := repository.NewSession(kvStore)

	// act
	err := sut.Create(t.Context(), session)

	// assert
	require.NoError(t, err)
	atExists, err := sut.AccessTokenExists(t.Context(), session.UserID, session.Tokens.Access.Value)
	require.NoError(t, err)
	assert.True(t, atExists)
	rtExists, err := sut.RefreshTokenExists(t.Context(), session.UserID, session.Tokens.Refresh.Value)
	require.NoError(t, err)
	assert.True(t, rtExists)
}

func TestSession_Delete(t *testing.T) {
	t.Parallel()

	// arrange
	kvStore := itest.NewKVS(t)
	session := model.MustNewSession(ulid.New())
	sut := repository.NewSession(kvStore)
	require.NoError(t, sut.Create(t.Context(), session))

	// act
	err := sut.Delete(t.Context(), session)

	// assert
	require.NoError(t, err)
	atExists, err := sut.AccessTokenExists(t.Context(), session.UserID, session.Tokens.Access.Value)
	require.NoError(t, err)
	assert.False(t, atExists)
	rtExists, err := sut.RefreshTokenExists(t.Context(), session.UserID, session.Tokens.Refresh.Value)
	require.NoError(t, err)
	assert.False(t, rtExists)
}

func TestSession_AccessTokenExists(t *testing.T) {
	t.Parallel()

	// arrange
	kvStore := itest.NewKVS(t)
	sut := repository.NewSession(kvStore)

	// act
	exists, err := sut.AccessTokenExists(t.Context(), "nonexistent", "nonexistent")

	// assert
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestSession_RefreshTokenExists(t *testing.T) {
	t.Parallel()

	// arrange
	kvStore := itest.NewKVS(t)
	sut := repository.NewSession(kvStore)

	// act
	exists, err := sut.RefreshTokenExists(t.Context(), "nonexistent", "nonexistent")

	// assert
	require.NoError(t, err)
	assert.False(t, exists)
}
