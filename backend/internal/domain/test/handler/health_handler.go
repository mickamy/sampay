package handler

import (
	"context"

	"connectrpc.com/connect"

	testv1 "github.com/mickamy/sampay/gen/test/v1"
	"github.com/mickamy/sampay/gen/test/v1/testv1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/lib/logger"
)

var _ testv1connect.HealthServiceHandler = (*HealthHandler)(nil)

type HealthHandler struct {
	infra *di.Infra
}

func NewHealthHandler(infra *di.Infra) *HealthHandler {
	return &HealthHandler{infra: infra}
}

func (h *HealthHandler) Check(
	ctx context.Context,
	r *connect.Request[testv1.CheckRequest],
) (*connect.Response[testv1.CheckResponse], error) {
	if _, err := h.infra.DB.ExecContext(ctx, "SELECT 1"); err != nil {
		logger.Error(ctx, "failed to ping database", "err", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if _, err := h.infra.WriterDB.ExecContext(ctx, "SELECT 1"); err != nil {
		logger.Error(ctx, "failed to ping writer database", "err", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if _, err := h.infra.ReaderDB.ExecContext(ctx, "SELECT 1"); err != nil {
		logger.Error(ctx, "failed to ping reader database", "err", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := h.infra.KVS.Ping(ctx); err != nil {
		logger.Error(ctx, "failed to ping KVS", "err", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&testv1.CheckResponse{}), nil
}
