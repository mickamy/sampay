package usecase_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/ulid"
)

func TestRefreshSession_Do(t *testing.T) {
	t.Parallel()

	userID := ulid.New()
	validSession := either.Must(authModel.NewSession(userID))

	tcs := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, kvs *kvs.KVS)
		input   func(t *testing.T) usecase.RefreshSessionInput
		assert  func(t *testing.T, got usecase.RefreshSessionOutput, err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, kvs *kvs.KVS) {
				require.NoError(t, authRepository.NewSession(kvs).Create(ctx, validSession))
			},
			input: func(t *testing.T) usecase.RefreshSessionInput {
				return usecase.RefreshSessionInput{RefreshToken: validSession.Tokens.Refresh.Value}
			},
			assert: func(t *testing.T, got usecase.RefreshSessionOutput, err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Tokens.Access.Value)
				assert.NotEmpty(t, got.Tokens.Refresh.Value)
			},
		},
		{
			name: "refresh token not found",
			arrange: func(t *testing.T, ctx context.Context, kvs *kvs.KVS) {
			},
			input: func(t *testing.T) usecase.RefreshSessionInput {
				return usecase.RefreshSessionInput{RefreshToken: validSession.Tokens.Refresh.Value}
			},
			assert: func(t *testing.T, got usecase.RefreshSessionOutput, err error) {
				require.ErrorIs(t, err, usecase.ErrRefreshSessionTokenNotFound)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctx := context.Background()
			kvStore := newKVS(t)
			tc.arrange(t, ctx, kvStore)

			// act
			got, err := usecase.NewRefreshSession(authRepository.NewSession(kvStore)).Do(ctx, tc.input(t))

			// assert
			tc.assert(t, got, err)
		})
	}
}
