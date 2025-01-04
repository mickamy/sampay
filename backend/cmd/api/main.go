package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	"connectrpc.com/connect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	s := newServer()

	fmt.Println("listening on port 8080...")
	if err := s.ListenAndServe(); err != nil {
		fmt.Println("failed to start server:", err)
		os.Exit(1)
	}
}

func newServer() *http.Server {
	api := http.NewServeMux()
	api.Handle(authv1connect.NewSessionServiceHandler(&SessionService{}))

	// TODO: restrict CORS
	corsHandler := cors.AllowAll().Handler(api)

	return &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(corsHandler, &http2.Server{}),
	}
}

type SessionService struct{}

func (s *SessionService) SignIn(ctx context.Context, req *connect.Request[authv1.SignInRequest]) (*connect.Response[authv1.SignInResponse], error) {
	panic("implement me")
}

func (s *SessionService) Refresh(ctx context.Context, req *connect.Request[authv1.RefreshRequest]) (*connect.Response[authv1.RefreshResponse], error) {
	panic("implement	me")
}

func (s *SessionService) SignOut(ctx context.Context, req *connect.Request[authv1.SignOutRequest]) (*connect.Response[authv1.SignOutResponse], error) {
	panic("implement me")
}

var _ authv1connect.SessionServiceHandler = (*SessionService)(nil)
