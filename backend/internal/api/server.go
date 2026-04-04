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
	"github.com/mickamy/sampay/internal/domain/auth"
	"github.com/mickamy/sampay/internal/domain/event"
	"github.com/mickamy/sampay/internal/domain/storage"
	"github.com/mickamy/sampay/internal/domain/test"
	"github.com/mickamy/sampay/internal/domain/user"
	"github.com/mickamy/sampay/internal/lib/logger"
)

func NewServer(infra *di.Infra) http.Server {
	interceptors := connect.WithInterceptors(interceptor.NewInterceptors(infra)...)

	api := http.NewServeMux()

	for _, route := range []router.Route{auth.Route, event.Route, storage.Route, test.Route, user.Route} {
		route(api, infra, interceptors)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if _, err := infra.WriterDB.ExecContext(r.Context(), "SELECT 1"); err != nil {
			logger.Error(r.Context(), "failed to ping writer database", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := infra.ReaderDB.ExecContext(r.Context(), "SELECT 1"); err != nil {
			logger.Error(r.Context(), "failed to ping reader database", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := infra.KVS.Ping(r.Context()); err != nil {
			logger.Error(r.Context(), "failed to ping KVS", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("/api/", http.StripPrefix("/api", api))

	return http.Server{
		Addr:              ":" + strconv.Itoa(config.API().Port),
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		ReadHeaderTimeout: 10 * time.Second,
	}
}
