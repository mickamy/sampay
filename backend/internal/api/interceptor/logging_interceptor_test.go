package interceptor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"

	testv1 "github.com/mickamy/sampay/gen/test/v1"
	"github.com/mickamy/sampay/gen/test/v1/testv1connect"
	"github.com/mickamy/sampay/internal/api/interceptor"
	"github.com/mickamy/sampay/internal/domain/test/handler"
)

func TestLoggingInterceptor(t *testing.T) {
	t.Parallel()

	// arrange
	test := func(ctx context.Context, req *connect.Request[testv1.TestRequest]) {
	}
	mux := http.NewServeMux()
	sut := interceptor.Logging()
	interceptors := connect.WithInterceptors(sut)
	mux.Handle(testv1connect.NewTestServiceHandler(&handler.TestHandler{Exec: test}, interceptors))
	server := httptest.NewServer(mux)
	defer server.Close()

	// act
	client := testv1connect.NewTestServiceClient(http.DefaultClient, server.URL)
	_, err := client.Test(t.Context(), connect.NewRequest(&testv1.TestRequest{}))

	// assert
	require.NoError(t, err)
	t.Log("Check the logs for the logging output")
}
