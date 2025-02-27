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

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/api/interceptor"
	"mickamy.com/sampay/internal/di"
	authRouter "mickamy.com/sampay/internal/domain/auth/router"
	commonRouter "mickamy.com/sampay/internal/domain/common/router"
	messageRouter "mickamy.com/sampay/internal/domain/message/router"
	notificationRouter "mickamy.com/sampay/internal/domain/notification/router"
	oauthRouter "mickamy.com/sampay/internal/domain/oauth/router"
	registrationRouter "mickamy.com/sampay/internal/domain/registration/router"
	userRouter "mickamy.com/sampay/internal/domain/user/router"
)

func NewServer(infras di.Infras) http.Server {
	authUCs := di.InitAuthUseCases(infras.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS)
	interceptors := connect.WithInterceptors(
		interceptor.Logging(),
		interceptor.I18N(),
		interceptor.Recovery(),
		interceptor.Authenticate(authUCs.AuthenticateUser, authUCs.AuthenticateAnonymousUser),
		interceptor.Cookie(),
	)

	api := http.NewServeMux()
	for _, route := range []func(mux *http.ServeMux, infras di.Infras, options ...connect.HandlerOption){
		authRouter.Route,
		commonRouter.Route,
		messageRouter.Route,
		notificationRouter.Route,
		oauthRouter.Route,
		registrationRouter.Route,
		userRouter.Route,
	} {
		route(api, infras, interceptors)
	}

	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api))
	corsHandler := cors.AllowAll().Handler(mux)

	return http.Server{
		Addr:    fmt.Sprintf(":%d", config.Common().Port),
		Handler: h2c.NewHandler(corsHandler, &http2.Server{}),
	}
}
