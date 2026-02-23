package auth

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mickamy/sampay/gen/auth/v1/authv1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/auth/handler"
)

func Route(mux *http.ServeMux, infra *di.Infra, options ...connect.HandlerOption) {
	mux.Handle(authv1connect.NewOAuthServiceHandler(handler.NewOAuth(infra), options...))
	mux.Handle(authv1connect.NewSessionServiceHandler(handler.NewSession(infra), options...))
}
