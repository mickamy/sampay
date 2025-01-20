package interceptor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/test/v1/testv1connect"
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	testv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/test/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/api/interceptor"
	"mickamy.com/sampay/internal/domain/test/handler"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/misc/i18n"
)

func TestRecovery(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name   string
		err    error
		assert func(t *testing.T, got *connect.Response[testv1.TestResponse], err error)
	}{
		{name: "no error", err: nil, assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
			assert.NoError(t, err)
		}},
		{name: "error", err: assert.AnError, assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
			require.Error(t, err)
			assert.Equalf(t, connect.CodeInternal, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
			connErr := new(connect.Error)
			require.ErrorAs(t, err, &connErr)
			require.Len(t, connErr.Details(), 1)
			detail := either.Must(connErr.Details()[0].Value())
			if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
				require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.CommonHandlerErrorInternal}), errMsg.Message)
			} else {
				require.Failf(t, "unexpected detail type", "got=%T", detail)
			}
		}},
	}

	for _, tc := range tsc {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			test := func(ctx context.Context, req *connect.Request[testv1.TestRequest]) {
				if err := tc.err; err != nil {
					panic(err)
				}
			}
			mux := http.NewServeMux()
			sut := interceptor.Recovery()
			interceptors := connect.WithInterceptors(interceptor.I18N(), sut)
			mux.Handle(testv1connect.NewTestServiceHandler(&handler.TestHandler{Exec: test}, interceptors))
			server := httptest.NewServer(mux)
			defer server.Close()

			// act
			client := testv1connect.NewTestServiceClient(http.DefaultClient, server.URL)
			got, err := client.Test(context.Background(), connect.NewRequest(&testv1.TestRequest{}))

			// assert
			tc.assert(t, got, err)
		})
	}
}
