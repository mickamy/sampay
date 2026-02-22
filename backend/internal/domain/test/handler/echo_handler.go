package handler

import (
	"context"

	"connectrpc.com/connect"

	testv1 "github.com/mickamy/sampay/gen/test/v1"
	"github.com/mickamy/sampay/gen/test/v1/testv1connect"
)

var _ testv1connect.EchoServiceHandler = (*EchoHandler)(nil)

type EchoHandler struct{}

func (h *EchoHandler) Echo(
	ctx context.Context,
	r *connect.Request[testv1.EchoRequest],
) (*connect.Response[testv1.EchoResponse], error) {
	res := &testv1.EchoResponse{Message: r.Msg.GetMessage()}
	return connect.NewResponse(res), nil
}
