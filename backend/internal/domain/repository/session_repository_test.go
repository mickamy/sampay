package repository_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/domain/model"
	"mickamy.com/sampay/internal/domain/repository"
	"mickamy.com/sampay/internal/lib/ulid"
	"mickamy.com/sampay/test/infra"
)

func TestSessionRepository_Create(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	userID := ulid.New()
	session, err := model.NewSession(userID)
	require.NoError(t, err)

	// act
	sut := repository.NewSession(infra.NewKVS(t))
	err = sut.Create(ctx, session)

	// assert
	require.NoError(t, err)
	exists, err := sut.AccessTokenExists(ctx, userID, session.Tokens.Access.Value)
	require.NoError(t, err)
	assert.True(t, exists)
	exists, err = sut.RefreshTokenExists(ctx, userID, session.Tokens.Refresh.Value)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestSessionRepository_Delete(t *testing.T) {
	t.Parallel()

	// arrange
	ctx := context.Background()
	userID := ulid.New()
	session, err := model.NewSession(userID)
	require.NoError(t, err)

	// act
	sut := repository.NewSession(infra.NewKVS(t))
	err = sut.Create(ctx, session)
	require.NoError(t, err)
	err = sut.Delete(ctx, session)

	// assert
	require.NoError(t, err)
	exists, err := sut.AccessTokenExists(ctx, userID, session.Tokens.Access.Value)
	require.NoError(t, err)
	assert.False(t, exists)
	exists, err = sut.RefreshTokenExists(ctx, userID, session.Tokens.Refresh.Value)
	require.NoError(t, err)
}
