package router

import (
	"net/http"

	"connectrpc.com/connect"

	"github.com/mickamy/sampay/internal/di"
)

type Route func(mux *http.ServeMux, infra *di.Infra, options ...connect.HandlerOption)
