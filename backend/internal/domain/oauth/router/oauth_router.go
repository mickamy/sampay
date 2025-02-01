package router

import (
	"net/http"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/oauth/v1/oauthv1connect"
	"connectrpc.com/connect"

	"mickamy.com/sampay/internal/di"
)

func Route(mux *http.ServeMux, infras di.Infras, options ...connect.HandlerOption) {
	handlers := di.InitOAuthHandlers(infras.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS)
	mux.Handle(oauthv1connect.NewOAuthServiceHandler(handlers.OAuth, options...))
}
