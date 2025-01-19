package router

import (
	"net/http"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/auth/v1/authv1connect"
	"connectrpc.com/connect"

	"mickamy.com/sampay/internal/di"
)

func Route(mux *http.ServeMux, infras di.Infras, options ...connect.HandlerOption) {
	handlers := di.InitAuthHandlers(infras.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS)
	mux.Handle(
		authv1connect.NewSessionServiceHandler(
			handlers.Session,
			options...,
		),
	)
}
