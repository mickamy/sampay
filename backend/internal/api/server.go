package api

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/mickamy/slogger"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"mickamy.com/sampay/internal/api/interceptor"
	"mickamy.com/sampay/internal/di"
	authRouter "mickamy.com/sampay/internal/domain/auth/router"
)

func NewServer(infras di.Infras) http.Server {
	api := http.NewServeMux()

	interceptors := connect.WithInterceptors(
		interceptor.Logging(),
		interceptor.I18N(),
	)

	for _, route := range []func(mux *http.ServeMux, infras di.Infras, options ...connect.HandlerOption){
		authRouter.Route,
	} {
		route(api, infras, interceptors)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := infras.Writer.Exec("SELECT 1").Error; err != nil {
			slogger.ErrorCtx(r.Context(), "failed to ping database", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := infras.Reader.Exec("SELECT 1").Error; err != nil {
			slogger.ErrorCtx(r.Context(), "failed to ping database", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := infras.KVS.Ping(context.Background()).Err(); err != nil {
			slogger.ErrorCtx(r.Context(), "failed to ping kvs", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("/api/", http.StripPrefix("/api", api))

	corsHandler := cors.AllowAll().Handler(mux)

	return http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: h2c.NewHandler(corsHandler, &http2.Server{}),
	}
}
