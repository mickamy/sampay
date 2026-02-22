package handler

import (
	"context"

	"connectrpc.com/connect"

	testv1 "github.com/mickamy/sampay/gen/test/v1"
	"github.com/mickamy/sampay/gen/test/v1/testv1connect"
)

var _ testv1connect.HealthServiceHandler = (*HealthHandler)(nil)

type HealthHandler struct{}

func (h *HealthHandler) Check(
	ctx context.Context,
	r *connect.Request[testv1.CheckRequest],
) (*connect.Response[testv1.CheckResponse], error) {
	return connect.NewResponse(&testv1.CheckResponse{}), nil
}
