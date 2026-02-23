package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/domain/auth/repository"
	"github.com/mickamy/sampay/internal/domain/auth/usecase"
	"github.com/mickamy/sampay/internal/infra/storage/kvs"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

func TestLogout_Do(t *testing.T) {
	t.Parallel()

	userID := ulid.New()
	validSession := model.MustNewSession(userID)
	anotherSession := model.MustNewSession(ulid.New())

	tests := []struct {
		name    string
		arrange func(t *testing.T, kvs *kvs.KVS)
		input   usecase.LogoutInput
		assert  func(t *testing.T, kvs *kvs.KVS, err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, kvs *kvs.KVS) {
				require.NoError(t, repository.NewSession(kvs).Create(t.Context(), validSession))
			},
			input: usecase.LogoutInput{
				AccessToken:  validSession.Tokens.Access.Value,
				RefreshToken: validSession.Tokens.Refresh.Value,
			},
			assert: func(t *testing.T, kvs *kvs.KVS, err error) {
				require.NoError(t, err)
				repo := repository.NewSession(kvs)
				exists, _ := repo.AccessTokenExists(t.Context(), userID, validSession.Tokens.Access.Value)
				assert.False(t, exists)
				exists, _ = repo.RefreshTokenExists(t.Context(), userID, validSession.Tokens.Refresh.Value)
				assert.False(t, exists)
			},
		},
		{
			name:    "invalid access token",
			arrange: func(t *testing.T, kvs *kvs.KVS) {},
			input: usecase.LogoutInput{
				AccessToken:  "invalid",
				RefreshToken: validSession.Tokens.Refresh.Value,
			},
			assert: func(t *testing.T, kvs *kvs.KVS, err error) {
				require.ErrorIs(t, err, usecase.ErrLogoutInvalidAccessToken)
			},
		},
		{
			name:    "invalid refresh token",
			arrange: func(t *testing.T, kvs *kvs.KVS) {},
			input: usecase.LogoutInput{
				AccessToken:  validSession.Tokens.Access.Value,
				RefreshToken: "invalid",
			},
			assert: func(t *testing.T, kvs *kvs.KVS, err error) {
				require.ErrorIs(t, err, usecase.ErrLogoutInvalidRefreshToken)
			},
		},
		{
			name:    "token mismatch",
			arrange: func(t *testing.T, kvs *kvs.KVS) {},
			input: usecase.LogoutInput{
				AccessToken:  validSession.Tokens.Access.Value,
				RefreshToken: anotherSession.Tokens.Refresh.Value,
			},
			assert: func(t *testing.T, kvs *kvs.KVS, err error) {
				require.ErrorIs(t, err, usecase.ErrLogoutTokenMismatch)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra := newInfra(t)
			tt.arrange(t, infra.KVS)

			sut := usecase.NewLogout(infra)
			_, err := sut.Do(t.Context(), tt.input)

			tt.assert(t, infra.KVS, err)
		})
	}
}
