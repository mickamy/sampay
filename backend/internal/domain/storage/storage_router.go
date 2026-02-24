package storage

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mickamy/sampay/gen/storage/v1/storagev1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/storage/handler"
)

func Route(mux *http.ServeMux, infra *di.Infra, options ...connect.HandlerOption) {
	mux.Handle(storagev1connect.NewStorageServiceHandler(handler.NewStorage(infra), options...))
}
