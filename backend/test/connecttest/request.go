package connecttest

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"connectrpc.com/connect"

	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
)

func NewRequest[T any](t *testing.T, ctx context.Context, message *T, header http.Header) *connect.Request[T] {
	t.Helper()

	req := connect.Request[T]{Msg: message}
	for k, vs := range header {
		for _, v := range vs {
			req.Header().Add(k, v)
		}
	}
	return &req
}

func NewAuthenticatedRequest[T any](t *testing.T, ctx context.Context, message *T, header http.Header, session authModel.Session, kvs *kvs.KVS) *connect.Request[T] {
	t.Helper()

	req := NewRequest(t, ctx, message, header)
	req.Header().Add("Authorization", "Bearer "+session.Tokens.Access.Value)
	if err := authRepository.NewSession(kvs).Create(ctx, session); err != nil {
		t.Fatal(fmt.Errorf("failed to persist session: %w", err))
	}
	return req
}

func NewAnonymousRequest[T any](t *testing.T, ctx context.Context, message *T, header http.Header, verification authModel.EmailVerification) *connect.Request[T] {
	t.Helper()

	req := NewRequest(t, ctx, message, header)
	req.Header().Add("Authorization", "Bearer "+verification.Verified.Token)
	return req
}
