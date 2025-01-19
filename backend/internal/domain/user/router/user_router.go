package router

import (
	"net/http"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/user/v1/userv1connect"
	"github.com/bufbuild/connect-go"

	"mickamy.com/sampay/internal/di"
)

func Route(mux *http.ServeMux, infras di.Infras, options ...connect.HandlerOption) {
	handlers := di.InitUserHandler(infras.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS)
	mux.Handle(userv1connect.NewUserServiceHandler(handlers.User, options...))
	mux.Handle(userv1connect.NewUserLinkServiceHandler(handlers.UserLink, options...))
	mux.Handle(userv1connect.NewUserProfileServiceHandler(handlers.UserProfile, options...))
}
