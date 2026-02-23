package interceptor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/mickamy/errx/cerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	ausecase "github.com/mickamy/sampay/internal/domain/auth/usecase"
	thandler "github.com/mickamy/sampay/internal/domain/test/handler"
	"github.com/mickamy/sampay/internal/test/itest"

	testv1 "github.com/mickamy/sampay/gen/test/v1"
	"github.com/mickamy/sampay/gen/test/v1/testv1connect"
	"github.com/mickamy/sampay/internal/api/interceptor"
	"github.com/mickamy/sampay/internal/di"
	amodel "github.com/mickamy/sampay/internal/domain/auth/model"
	arepository "github.com/mickamy/sampay/internal/domain/auth/repository"
	"github.com/mickamy/sampay/internal/infra/storage/kvs"
	"github.com/mickamy/sampay/internal/lib/either"
	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/misc/i18n"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
	"github.com/mickamy/sampay/internal/test/ctest"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	validSession := either.Must(amodel.NewSession(uuid.NewString()))

	tests := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, kvs *kvs.KVS) string
		assert  func(t *testing.T, got *connect.Response[testv1.TestResponse], err error)
		want    string
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, kvs *kvs.KVS) string {
				require.NoError(t, arepository.NewSession(kvs).Create(ctx, validSession))
				return "Bearer " + validSession.Tokens.Access.Value
			},
			assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
				require.NoError(t, err)
			},
			want: validSession.UserID,
		},
		{
			name: "fail (token not set)",
			arrange: func(t *testing.T, ctx context.Context, kvs *kvs.KVS) string {
				require.NoError(t, arepository.NewSession(kvs).Create(ctx, validSession))
				return ""
			},
			assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
				require.Error(t, err)
				connErr := ctest.ConnErr(t, err)
				t.Logf("connErr: %+v", connErr)
				ctest.AssertCode(t, connect.CodeUnauthenticated, connErr)
				localized := ctest.LocalizedMessage(t, connErr)
				assert.Equal(t, i18n.Japanese(messages.AuthUseCaseErrorSessionNotSet()), localized)
				fieldViolation := ctest.FieldViolation(t, connErr)
				assert.Equal(t, "access_token", fieldViolation.GetField())
				assert.Equal(t, i18n.Japanese(messages.AuthUseCaseErrorSessionNotSet()), fieldViolation.GetDescription())
			},
		},
		{
			name: "fail (invalid token set)",
			arrange: func(t *testing.T, ctx context.Context, kvs *kvs.KVS) string {
				require.NoError(t, arepository.NewSession(kvs).Create(ctx, validSession))
				return "Bearer " + validSession.Tokens.Access.Value + "invalid"
			},
			assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
				require.Error(t, err)
				connErr := ctest.ConnErr(t, err)
				ctest.AssertCode(t, connect.CodeUnauthenticated, connErr)
				localized := ctest.LocalizedMessage(t, connErr)
				assert.Equal(t, i18n.Japanese(messages.AuthUseCaseErrorTokenInvalid()), localized)
				fieldViolation := ctest.FieldViolation(t, connErr)
				assert.Equal(t, "access_token", fieldViolation.GetField())
				assert.Equal(t, i18n.Japanese(messages.AuthUseCaseErrorTokenInvalid()), fieldViolation.GetDescription())
			},
		},
		{
			name: "fail (refresh token set)",
			arrange: func(t *testing.T, ctx context.Context, kvs *kvs.KVS) string {
				require.NoError(t, arepository.NewSession(kvs).Create(ctx, validSession))
				return "Bearer " + validSession.Tokens.Refresh.Value
			},
			assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
				require.Error(t, err)
				connErr := ctest.ConnErr(t, err)
				ctest.AssertCode(t, connect.CodeUnauthenticated, connErr)
				localized := ctest.LocalizedMessage(t, connErr)
				assert.Equal(t, i18n.Japanese(messages.AuthUseCaseErrorSessionNotFound()), localized)
				fieldViolation := ctest.FieldViolation(t, connErr)
				assert.Equal(t, "access_token", fieldViolation.GetField())
				assert.Equal(t, i18n.Japanese(messages.AuthUseCaseErrorSessionNotFound()), fieldViolation.GetDescription())
			},
		},
		{
			name: "fail (session not exists)",
			arrange: func(t *testing.T, ctx context.Context, kvs *kvs.KVS) string {
				return "Bearer " + validSession.Tokens.Access.Value
			},
			assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
				require.Error(t, err)
				connErr := ctest.ConnErr(t, err)
				ctest.AssertCode(t, connect.CodeUnauthenticated, connErr)
				localized := ctest.LocalizedMessage(t, connErr)
				assert.Equal(t, i18n.Japanese(messages.AuthUseCaseErrorSessionNotFound()), localized)
				fieldViolation := ctest.FieldViolation(t, connErr)
				assert.Equal(t, "access_token", fieldViolation.GetField())
				assert.Equal(t, i18n.Japanese(messages.AuthUseCaseErrorSessionNotFound()), fieldViolation.GetDescription())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctx := t.Context()
			kvStore := itest.NewKVS(t)
			authorization := tt.arrange(t, ctx, kvStore)
			test := func(ctx context.Context, req *connect.Request[testv1.TestRequest]) {
				userID, err := contexts.AuthenticatedUserID(ctx)
				logger.Info(ctx, "test", "user_id", userID, "name", t.Name())
				assert.Equal(t, tt.want, userID)
				if tt.want == "" {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			}
			mux := http.NewServeMux()
			uc := ausecase.NewAuthenticate(&di.Infra{KVS: kvStore})
			sut := interceptor.Authenticate(uc)
			interceptors := connect.WithInterceptors(cerr.NewInterceptor(
				cerr.WithLocaleFunc(func(header http.Header) string { return header.Get("Accept-Language") }),
				cerr.WithDefaultLocale(language.Japanese),
			), interceptor.I18N(), sut)
			mux.Handle(testv1connect.NewTestServiceHandler(&thandler.TestHandler{Exec: test}, interceptors))
			server := httptest.NewServer(mux)
			defer server.Close()

			// act
			client := testv1connect.NewTestServiceClient(http.DefaultClient, server.URL)
			req := connect.NewRequest(&testv1.TestRequest{})
			req.Header().Add("Authorization", authorization)
			got, err := client.Test(ctx, req)

			// assert
			tt.assert(t, got, err)
		})
	}
}
