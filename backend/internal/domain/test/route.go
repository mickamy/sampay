package test

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mickamy/sampay/gen/test/v1/testv1connect"
	"github.com/mickamy/sampay/internal/domain/test/handler"
)

func Route(mux *http.ServeMux, options ...connect.HandlerOption) {
	mux.Handle(testv1connect.NewHealthServiceHandler(&handler.HealthHandler{}, options...))
}
