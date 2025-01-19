package router

import (
	"net/http"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/registration/v1/registrationv1connect"
	"github.com/bufbuild/connect-go"

	"mickamy.com/sampay/internal/di"
)

func Route(mux *http.ServeMux, infras di.Infras, options ...connect.HandlerOption) {
	handlers := di.InitRegistrationHandlers(infras.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS)
	mux.Handle(registrationv1connect.NewAccountServiceHandler(handlers.Account, options...))
	mux.Handle(registrationv1connect.NewEmailVerificationServiceHandler(handlers.EmailVerification, options...))
	mux.Handle(registrationv1connect.NewOnboardingServiceHandler(handlers.Onboarding, options...))
	mux.Handle(registrationv1connect.NewUsageCategoryServiceHandler(handlers.UsageCategory, options...))
}
