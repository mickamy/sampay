package router

import (
	"net/http"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/message/v1/messagev1connect"
	"connectrpc.com/connect"

	"mickamy.com/sampay/internal/di"
)

func Route(mux *http.ServeMux, infras di.Infras, options ...connect.HandlerOption) {
	handlers := di.InitMessageHandlers(infras.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS)
	mux.Handle(messagev1connect.NewMessageServiceHandler(handlers.Message, options...))
}
