package test

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mickamy/sampay/gen/test/v1/testv1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/test/handler"
)

func Route(mux *http.ServeMux, infra *di.Infra, options ...connect.HandlerOption) {
	mux.Handle(testv1connect.NewHealthServiceHandler(handler.NewHealthHandler(infra), options...))
}
