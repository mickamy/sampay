package router

import (
	"net/http"

	"connectrpc.com/connect"
)

type Route func(mux *http.ServeMux, options ...connect.HandlerOption)
