package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/domain/auth/repository"
	"github.com/mickamy/sampay/internal/domain/auth/usecase"
	"github.com/mickamy/sampay/internal/infra/storage/kvs"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

func TestRefreshToken_Do(t *testing.T) {
	t.Parallel()

	userID := ulid.New()
	validSession := model.MustNewSession(userID)

	tests := []struct {
		name    string
		arrange func(t *testing.T, kvs *kvs.KVS)
		token   string
		assert  func(t *testing.T, got usecase.RefreshTokenOutput, err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, kvs *kvs.KVS) {
				require.NoError(t, repository.NewSession(kvs).Create(t.Context(), validSession))
			},
			token: validSession.Tokens.Refresh.Value,
			assert: func(t *testing.T, got usecase.RefreshTokenOutput, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, got.Tokens.Access.Value)
				require.NotEmpty(t, got.Tokens.Refresh.Value)
			},
		},
		{
			name:    "token not set",
			arrange: func(t *testing.T, kvs *kvs.KVS) {},
			token:   "",
			assert: func(t *testing.T, got usecase.RefreshTokenOutput, err error) {
				require.ErrorIs(t, err, usecase.ErrRefreshTokenNotSet)
			},
		},
		{
			name:    "refresh token not found",
			arrange: func(t *testing.T, kvs *kvs.KVS) {},
			token:   validSession.Tokens.Refresh.Value,
			assert: func(t *testing.T, got usecase.RefreshTokenOutput, err error) {
				require.ErrorIs(t, err, usecase.ErrRefreshTokenNotFound)
			},
		},
		{
			name: "invalid refresh token",
			arrange: func(t *testing.T, kvs *kvs.KVS) {
				require.NoError(t, repository.NewSession(kvs).Create(t.Context(), validSession))
			},
			token: validSession.Tokens.Refresh.Value + "_invalid",
			assert: func(t *testing.T, got usecase.RefreshTokenOutput, err error) {
				require.ErrorIs(t, err, usecase.ErrRefreshTokenInvalid)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			infra := newInfra(t)
			tt.arrange(t, infra.KVS)

			sut := usecase.NewRefreshToken(infra)
			got, err := sut.Do(t.Context(), usecase.RefreshTokenInput{
				Token: tt.token,
			})

			tt.assert(t, got, err)
		})
	}
}
