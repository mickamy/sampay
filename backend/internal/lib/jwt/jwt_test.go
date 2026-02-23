package jwt_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/lib/jwt"
)

func TestJWT_New(t *testing.T) {
	t.Parallel()

	// arrange
	id := uuid.NewString()

	// act
	token, err := jwt.New(id)

	// assert
	require.NoError(t, err)
	assert.NotEmpty(t, token.Access.Value)
	assert.NotEmpty(t, token.Access.ExpiresAt)
	assert.NotEmpty(t, token.Refresh.Value)
	assert.NotEmpty(t, token.Refresh.ExpiresAt)
}

func TestJWT_Verify(t *testing.T) {
	t.Parallel()

	// arrange
	id := uuid.NewString()
	token, err := jwt.New(id)
	require.NoError(t, err)

	// act
	accessClaim, err := jwt.Verify(token.Access.Value)
	require.NoError(t, err)
	refreshClaim, err := jwt.Verify(token.Refresh.Value)
	require.NoError(t, err)

	// assert
	assert.Equal(t, accessClaim["id"], id)
	assert.NotEmpty(t, accessClaim["exp"])
	assert.Equal(t, refreshClaim["id"], id)
	assert.NotEmpty(t, refreshClaim["exp"], token.Access.ExpiresAt)
}

func TestJWT_ExtractID(t *testing.T) {
	t.Parallel()

	// arrange
	userID := uuid.NewString()
	token, err := jwt.New(userID)
	require.NoError(t, err)

	tests := []struct {
		name     string
		token    jwt.Token
		expected string
	}{
		{
			name:     "from access token",
			token:    token.Access,
			expected: userID,
		},
		{
			name:     "from refresh token",
			token:    token.Refresh,
			expected: userID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// act
			actual, err := jwt.ExtractID(tt.token.Value)

			// assert
			require.NoError(t, err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
