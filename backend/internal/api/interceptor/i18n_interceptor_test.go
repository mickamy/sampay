package interceptor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/test/v1/testv1connect"
	testv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/test/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/api/interceptor"
	"mickamy.com/sampay/internal/domain/test/handler"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/language"
)

func TestI18N(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name         string
		acceptLang   string
		expectedLang language.Type
	}{
		{name: "success", acceptLang: "ja", expectedLang: language.Japanese},
		{name: "success (not supported language)", acceptLang: "en", expectedLang: language.Japanese},
		{name: "success (empty header)", acceptLang: "", expectedLang: language.Japanese},
		{name: "success (invalid language format)", acceptLang: "test", expectedLang: language.Japanese},
	}

	for _, tc := range tsc {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			test := func(ctx context.Context, req *connect.Request[testv1.TestRequest]) {
				lang := contexts.MustLanguage(ctx)
				require.Equal(t, tc.expectedLang, lang)
			}
			mux := http.NewServeMux()
			sut := interceptor.I18N()
			interceptors := connect.WithInterceptors(sut)
			mux.Handle(testv1connect.NewTestServiceHandler(&handler.TestHandler{Exec: test}, interceptors))
			server := httptest.NewServer(mux)
			defer server.Close()

			// act
			client := testv1connect.NewTestServiceClient(http.DefaultClient, server.URL)
			_, err := client.Test(context.Background(), connect.NewRequest(&testv1.TestRequest{}))

			// assert
			require.NoError(t, err)
		})
	}
}
