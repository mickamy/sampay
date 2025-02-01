package router

import (
	"net/http"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/registration/v1/registrationv1connect"
	"connectrpc.com/connect"

	"mickamy.com/sampay/internal/di"
)

func Route(mux *http.ServeMux, infras di.Infras, options ...connect.HandlerOption) {
	handlers := di.InitRegistrationHandlers(infras.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS)
	mux.Handle(registrationv1connect.NewOnboardingServiceHandler(handlers.Onboarding, options...))
	mux.Handle(registrationv1connect.NewUsageCategoryServiceHandler(handlers.UsageCategory, options...))
}
