package interceptor_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/registration/v1/registrationv1connect"
	"buf.build/gen/go/mickamy/sampay/connectrpc/go/test/v1/testv1connect"
	registrationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/registration/v1"
	testv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/test/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/api/interceptor"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	testHandler "mickamy.com/sampay/internal/domain/test/handler"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	userRepository "mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/ulid"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	t.Run("authenticated", func(t *testing.T) {
		t.Parallel()

		validSession := either.Must(authModel.NewSession(ulid.New()))

		tsc := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string
			assert  func(t *testing.T, got *connect.Response[testv1.TestResponse], err error)
			want    string
		}{
			{
				name: "success",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string {
					u := userFixture.User(func(m *userModel.User) {
						m.ID = validSession.UserID
					})
					require.NoError(t, userRepository.NewUser(db).Create(ctx, &u))
					require.NoError(t, authRepository.NewSession(kvs).Create(ctx, validSession))
					return "Bearer " + validSession.Tokens.Access.Value
				},
				assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
					require.NoError(t, err)
				},
				want: validSession.UserID,
			},
			{
				name: "fail (token not set)",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string {
					u := userFixture.User(func(m *userModel.User) {
						m.ID = validSession.UserID
					})
					require.NoError(t, userRepository.NewUser(db).Create(ctx, &u))
					require.NoError(t, authRepository.NewSession(kvs).Create(ctx, validSession))
					return ""
				},
				assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
					require.Error(t, err)
					assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
				},
			},
			{
				name: "fail (invalid token set)",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string {
					u := userFixture.User(func(m *userModel.User) {
						m.ID = validSession.UserID
					})
					require.NoError(t, userRepository.NewUser(db).Create(ctx, &u))
					require.NoError(t, authRepository.NewSession(kvs).Create(ctx, validSession))
					return "Bearer " + validSession.Tokens.Access.Value + "invalid"
				},
				assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
					require.Error(t, err)
					assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
				},
			},
			{
				name: "fail (user not exists)",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string {
					require.NoError(t, authRepository.NewSession(kvs).Create(ctx, validSession))
					return "Bearer " + validSession.Tokens.Access.Value
				},
				assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
					require.Error(t, err)
					assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
				},
			},
			{
				name: "fail (refresh token set)",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string {
					u := userFixture.User(func(m *userModel.User) {
						m.ID = validSession.UserID
					})
					require.NoError(t, userRepository.NewUser(db).Create(ctx, &u))
					require.NoError(t, authRepository.NewSession(kvs).Create(ctx, validSession))
					return "Bearer " + validSession.Tokens.Refresh.Value
				},
				assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
					require.Error(t, err)
					assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
				},
			},
			{
				name: "fail (session not exists)",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string {
					u := userFixture.User(func(m *userModel.User) {
						m.ID = validSession.UserID
					})
					require.NoError(t, userRepository.NewUser(db).Create(ctx, &u))
					return "Bearer " + validSession.Tokens.Access.Value
				},
				assert: func(t *testing.T, got *connect.Response[testv1.TestResponse], err error) {
					require.Error(t, err)
					assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
				},
			},
		}

		for _, tc := range tsc {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				// arrange
				ctx := context.Background()
				db := newReadWriter(t)
				kvStore := newKVS(t)
				authorization := tc.arrange(t, ctx, db.WriterDB(), kvStore)
				test := func(ctx context.Context, req *connect.Request[testv1.TestRequest]) {
					authorizedUser, err := contexts.AuthenticatedUserID(ctx)
					slogger.InfoCtx(ctx, "authorizedUser", "authorizedUser", authorizedUser, "name", t.Name())
					assert.Equal(t, tc.want, authorizedUser)
					if tc.want == "" {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
					}
				}
				mux := http.NewServeMux()
				ucs := di.InitAuthUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), kvStore)
				sut := interceptor.Authenticate(ucs.AuthenticateUser, ucs.AuthenticateAnonymousUser)
				interceptors := connect.WithInterceptors(sut)
				mux.Handle(testv1connect.NewTestServiceHandler(&testHandler.TestHandler{Exec: test}, interceptors))
				server := httptest.NewServer(mux)
				defer server.Close()

				// act
				client := testv1connect.NewTestServiceClient(http.DefaultClient, server.URL)
				req := connect.NewRequest(&testv1.TestRequest{})
				req.Header().Add("Authorization", authorization)
				got, err := client.Test(ctx, req)

				// assert
				tc.assert(t, got, err)
			})
		}
	})

	t.Run("anonymous", func(t *testing.T) {
		t.Parallel()

		tsc := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string
			assert  func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error)
		}{

			{
				name: "success",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string {
					verification := authFixture.EmailVerificationVerified(nil)
					require.NoError(t, db.WithContext(ctx).Create(&verification).Error)
					return "Bearer " + verification.Verified.Token
				},
				assert: func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
					require.NoError(t, err)
				},
			},
			{
				name: "fail (token not set)",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string {
					verification := authFixture.EmailVerificationVerified(nil)
					require.NoError(t, db.WithContext(ctx).Create(&verification).Error)
					return ""
				},
				assert: func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
					require.Error(t, err)
					assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
				},
			},
			{
				name: "fail (invalid token set)",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string {
					verification := authFixture.EmailVerificationVerified(nil)
					require.NoError(t, db.WithContext(ctx).Create(&verification).Error)
					return "Bearer " + verification.Verified.Token + "invalid"
				},
				assert: func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
					require.Error(t, err)
					assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
				},
			},
			{
				name: "fail (verification not exists)",
				arrange: func(t *testing.T, ctx context.Context, db *database.DB, kvs *kvs.KVS) string {
					verification := authFixture.EmailVerificationVerified(nil)
					return "Bearer " + verification.Verified.Token
				},
				assert: func(t *testing.T, got *connect.Response[registrationv1.GetOnboardingStepResponse], err error) {
					require.Error(t, err)
					assert.Equal(t, connect.CodeUnauthenticated, connect.CodeOf(err))
				},
			},
		}

		for _, tc := range tsc {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				// arrange
				ctx := context.Background()
				db := newReadWriter(t)
				kvStore := newKVS(t)
				if tc.arrange == nil {
					require.Fail(t, "name="+tc.name)
				}
				authorization := tc.arrange(t, ctx, db.WriterDB(), kvStore)
				mux := http.NewServeMux()
				ucs := di.InitAuthUseCases(db.WriterDB(), db, db.Writer(), db.Reader(), kvStore)
				sut := interceptor.Authenticate(ucs.AuthenticateUser, ucs.AuthenticateAnonymousUser)
				interceptors := connect.WithInterceptors(sut)
				onboarding := di.InitRegistrationHandlers(db.WriterDB(), db, db.Writer(), db.Reader(), newKVS(t)).Onboarding
				mux.Handle(registrationv1connect.NewOnboardingServiceHandler(onboarding, interceptors))
				server := httptest.NewServer(mux)
				defer server.Close()

				// act
				client := registrationv1connect.NewOnboardingServiceClient(http.DefaultClient, server.URL)
				req := connect.NewRequest(&registrationv1.GetOnboardingStepRequest{})
				req.Header().Add("Authorization", authorization)
				got, err := client.GetOnboardingStep(ctx, req)

				// assert
				tc.assert(t, got, err)
			})
		}
	})
}
