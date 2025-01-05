package handler

import (
	"context"

	testv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/test/v1"
	"connectrpc.com/connect"
)

type EchoHandler struct{}

func (h *EchoHandler) Echo(
	ctx context.Context,
	req *connect.Request[testv1.EchoRequest],
) (*connect.Response[testv1.EchoResponse], error) {
	res := &testv1.EchoResponse{Message: req.Msg.Message}
	return connect.NewResponse(res), nil
}
