package router

import (
	"net/http"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/notification/v1/notificationv1connect"
	"connectrpc.com/connect"

	"mickamy.com/sampay/internal/di"
)

func Route(mux *http.ServeMux, infras di.Infras, options ...connect.HandlerOption) {
	handlers := di.InitNotificationHandlers(infras.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS)
	mux.Handle(notificationv1connect.NewNotificationServiceHandler(handlers.Notification, options...))
}
