package handler

import (
	"context"

	"connectrpc.com/connect"

	testv1 "github.com/mickamy/sampay/gen/test/v1"
	"github.com/mickamy/sampay/gen/test/v1/testv1connect"
)

var _ testv1connect.TestServiceHandler = (*TestHandler)(nil)

type TestHandler struct {
	Exec func(ctx context.Context, req *connect.Request[testv1.TestRequest])
}

func (h *TestHandler) Test(
	ctx context.Context,
	r *connect.Request[testv1.TestRequest],
) (*connect.Response[testv1.TestResponse], error) {
	if h.Exec != nil {
		h.Exec(ctx, r)
	}
	return connect.NewResponse(&testv1.TestResponse{}), nil
}
