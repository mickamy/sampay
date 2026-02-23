package interceptor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	testv1 "github.com/mickamy/sampay/gen/test/v1"
	"github.com/mickamy/sampay/gen/test/v1/testv1connect"
	"github.com/mickamy/sampay/internal/api/interceptor"
	"github.com/mickamy/sampay/internal/domain/test/handler"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

func TestI18N(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name         string
		acceptLang   string
		expectedLang language.Tag
	}{
		{name: "success", acceptLang: "ja", expectedLang: language.Japanese},
		{name: "success (not supported language)", acceptLang: "fr", expectedLang: language.Japanese}, // french
		{name: "success (empty header)", acceptLang: "", expectedLang: language.Japanese},
		{name: "success (invalid language format)", acceptLang: "test", expectedLang: language.Japanese},
	}

	for _, tc := range tsc {
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
			_, err := client.Test(t.Context(), connect.NewRequest(&testv1.TestRequest{}))

			// assert
			require.NoError(t, err)
		})
	}
}
