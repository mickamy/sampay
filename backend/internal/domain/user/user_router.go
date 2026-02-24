package user

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mickamy/sampay/gen/user/v1/userv1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/user/handler"
)

func Route(mux *http.ServeMux, infra *di.Infra, options ...connect.HandlerOption) {
	mux.Handle(userv1connect.NewPaymentMethodServiceHandler(handler.NewPaymentMethod(infra), options...))
}
