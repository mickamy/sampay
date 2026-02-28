package event

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mickamy/sampay/gen/event/v1/eventv1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/event/handler"
)

func Route(mux *http.ServeMux, infra *di.Infra, options ...connect.HandlerOption) {
	mux.Handle(eventv1connect.NewEventServiceHandler(handler.NewEventService(infra), options...))
	mux.Handle(eventv1connect.NewEventProfileServiceHandler(handler.NewEventProfile(infra), options...))
}
