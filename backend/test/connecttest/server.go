package connecttest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"

	"mickamy.com/sampay/internal/api/interceptor"
	"mickamy.com/sampay/internal/di"
)

func NewServer(t *testing.T, infras di.Infras, newService func([]connect.Interceptor) (string, http.Handler)) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()
	mux.Handle(
		newService(NewInterceptors(infras)),
	)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	return server
}

func NewInterceptors(infras di.Infras) []connect.Interceptor {
	return []connect.Interceptor{
		interceptor.Logging(),
		interceptor.I18N(),
		interceptor.Recovery(),
		interceptor.Authenticate(di.InitAuthUseCases(infras.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).AuthenticateUser),
		interceptor.Cookie(),
	}
}
