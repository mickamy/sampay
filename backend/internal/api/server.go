package api

import (
	"net/http"
	"strconv"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/api/interceptor"
	"github.com/mickamy/sampay/internal/api/router"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/test"
)

func NewServer(infra *di.Infra) http.Server {
	interceptors := connect.WithInterceptors(interceptor.NewInterceptors(infra)...)

	api := http.NewServeMux()

	for _, route := range []router.Route{test.Route} {
		route(api, infra, interceptors)
	}

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api))

	return http.Server{
		Addr:              ":" + strconv.Itoa(config.API().Port),
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 10 * time.Second,
	}
}
