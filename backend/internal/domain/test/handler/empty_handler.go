package handler

import (
	"context"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/test/v1/testv1connect"
	testv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/test/v1"
	"connectrpc.com/connect"
)

type TestHandler struct {
	Exec func(ctx context.Context, req *connect.Request[testv1.TestRequest])
}

func (h *TestHandler) Test(
	ctx context.Context,
	req *connect.Request[testv1.TestRequest],
) (*connect.Response[testv1.TestResponse], error) {
	if h.Exec != nil {
		h.Exec(ctx, req)
	}
	return connect.NewResponse(&testv1.TestResponse{}), nil
}

var (
	_ testv1connect.TestServiceHandler = (*TestHandler)(nil)
)
