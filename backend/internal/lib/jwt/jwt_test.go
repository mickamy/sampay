package jwt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"mickamy.com/sampay/internal/lib/jwt"
	"mickamy.com/sampay/internal/lib/ulid"
)

func TestJWT_New(t *testing.T) {
	t.Parallel()

	// arrange
	id := ulid.New()

	// act
	token, err := jwt.New(id)

	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token.Access.Value)
	assert.NotEmpty(t, token.Access.ExpiresAt)
	assert.NotEmpty(t, token.Refresh.Value)
	assert.NotEmpty(t, token.Refresh.ExpiresAt)
}

func TestJWT_Verify(t *testing.T) {
	t.Parallel()

	// arrange
	id := ulid.New()
	token, err := jwt.New(id)
	assert.NoError(t, err)

	// act
	accessClaim, err := jwt.Verify(token.Access.Value)
	assert.NoError(t, err)
	refreshClaim, err := jwt.Verify(token.Refresh.Value)
	assert.NoError(t, err)

	// assert
	assert.Equal(t, accessClaim["id"], id)
	assert.NotEmpty(t, accessClaim["exp"])
	assert.Equal(t, refreshClaim["id"], id)
	assert.Equal(t, refreshClaim["jwt"], token.Access.Value)
	assert.NotEmpty(t, refreshClaim["exp"], token.Access.ExpiresAt)
}

func TestJWT_ExtractID(t *testing.T) {
	t.Parallel()

	// arrange
	userID := ulid.New()
	token, err := jwt.New(userID)
	assert.NoError(t, err)

	tcs := []struct {
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
	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			c := c
			t.Parallel()

			// act
			actual, err := jwt.ExtractID(c.token.Value)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestJWT_IsRefreshTokenClaims(t *testing.T) {
	t.Parallel()

	// arrange
	userID := ulid.New()
	tokens, err := jwt.New(userID)
	assert.NoError(t, err)

	tcs := []struct {
		name  string
		token jwt.Token
		want  bool
	}{
		{
			name:  "access token",
			token: tokens.Access,
			want:  false,
		},
		{
			name:  "refresh token",
			token: tokens.Refresh,
			want:  true,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			claims, err := jwt.Verify(tc.token.Value)
			assert.NoError(t, err)

			// act
			got := jwt.IsRefreshTokenClaims(claims)

			// assert
			assert.Equal(t, tc.want, got)
		})
	}
}
