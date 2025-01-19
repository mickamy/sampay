package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/auth/v1/authv1connect"
	authv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/auth/v1"
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	"connectrpc.com/connect"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authFixture "mickamy.com/sampay/internal/domain/auth/fixture"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	authRepository "mickamy.com/sampay/internal/domain/auth/repository"
	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/ptr"
	"mickamy.com/sampay/internal/misc/i18n"
	"mickamy.com/sampay/test/connecttest"
)

func TestSession_SignIn(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.SignInRequest
		assert  func(t *testing.T, got *connect.Response[authv1.SignInResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.SignInRequest {
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)
				return &authv1.SignInRequest{
					Email:    auth.Identifier,
					Password: commonFixture.Password,
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.SignInResponse], err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Msg.UserId)
				assert.NotEmpty(t, got.Msg.Tokens.Access.Value)
				assert.NotEmpty(t, got.Msg.Tokens.Access.ExpiresAt)
				assert.NotEmpty(t, got.Msg.Tokens.Refresh.Value)
				assert.NotEmpty(t, got.Msg.Tokens.Refresh.ExpiresAt)
			},
		},
		{
			name: "fail (email not found)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.SignInRequest {
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)
				return &authv1.SignInRequest{
					Email:    gofakeit.GlobalFaker.Email(),
					Password: commonFixture.Password,
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.SignInResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
					expectedMsg := i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseErrorInvalid_email_password})
					assert.Equal(t, expectedMsg, errMsg.Message)
				} else {
					require.Failf(t, "unexpected detail type", "got=%T", detail)
				}
			},
		},
		{
			name: "fail (password not found)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.SignInRequest {
				auth := authFixture.AuthenticationEmailPassword(func(m *authModel.Authentication) {
					m.UserID = userID
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&auth).Error)
				return &authv1.SignInRequest{
					Email:    auth.Identifier,
					Password: commonFixture.Password + "invalid",
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.SignInResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
					expectedMsg := i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseErrorInvalid_email_password})
					assert.Equal(t, expectedMsg, errMsg.Message)
				} else {
					require.Failf(t, "unexpected detail type", "got=%T", detail)
				}
			},
		},
	}

	for _, tc := range tsc {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctx := context.Background()
			infras := di.NewInfras(newReadWriter(t), newKVS(t))
			user := userFixture.User(nil)
			require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
			ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
			req := tc.arrange(t, ctx, infras, user.ID)
			server := newSessionServer(t, infras)

			// act
			client := authv1connect.NewSessionServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.SignIn(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestSession_Refresh(t *testing.T) {
	t.Parallel()

	user := userFixture.User(nil)

	t.Run("in message", func(t *testing.T) {
		tsc := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.RefreshRequest
			assert  func(t *testing.T, got *connect.Response[authv1.RefreshResponse], err error)
		}{
			{
				name: "success",
				arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.RefreshRequest {
					session := authModel.MustNewSession(userID)
					require.NoError(t, authRepository.NewSession(infras.KVS).Create(ctx, session))
					return &authv1.RefreshRequest{
						RefreshToken: ptr.Of(session.Tokens.Refresh.Value),
					}
				},
				assert: func(t *testing.T, got *connect.Response[authv1.RefreshResponse], err error) {
					require.NoError(t, err)
					assert.NotEmpty(t, got.Msg.Tokens.Access.Value)
					assert.NotEmpty(t, got.Msg.Tokens.Access.ExpiresAt)
					assert.NotEmpty(t, got.Msg.Tokens.Refresh.Value)
					assert.NotEmpty(t, got.Msg.Tokens.Refresh.ExpiresAt)
				},
			},
			{
				name: "fail (refresh token not found)",
				arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.RefreshRequest {
					session := authModel.MustNewSession(userID)
					require.NoError(t, authRepository.NewSession(infras.KVS).Create(ctx, session))
					return &authv1.RefreshRequest{}
				},
				assert: func(t *testing.T, got *connect.Response[authv1.RefreshResponse], err error) {
					require.Error(t, err)
					assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
					connErr := new(connect.Error)
					require.ErrorAs(t, err, &connErr)
					require.Len(t, connErr.Details(), 1)
					detail := either.Must(connErr.Details()[0].Value())
					if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
						expectedMsg := i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseErrorInvalid_refresh_token})
						assert.Equal(t, expectedMsg, errMsg.Message)
					} else {
						require.Failf(t, "unexpected detail type", "got=%T", detail)
					}
				},
			},
			{
				name: "fail (invalid refresh token)",
				arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.RefreshRequest {
					session := authModel.MustNewSession(userID)
					require.NoError(t, authRepository.NewSession(infras.KVS).Create(ctx, session))
					return &authv1.RefreshRequest{
						RefreshToken: ptr.Of(session.Tokens.Refresh.Value + "invalid"),
					}
				},
				assert: func(t *testing.T, got *connect.Response[authv1.RefreshResponse], err error) {
					require.Error(t, err)
					assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
					connErr := new(connect.Error)
					require.ErrorAs(t, err, &connErr)
					require.Len(t, connErr.Details(), 1)
					detail := either.Must(connErr.Details()[0].Value())
					if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
						expectedMsg := i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseErrorInvalid_refresh_token})
						assert.Equal(t, expectedMsg, errMsg.Message)
					} else {
						require.Failf(t, "unexpected detail type", "got=%T", detail)
					}
				},
			},
		}

		for _, tc := range tsc {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				// arrange
				ctx := context.Background()
				infras := di.NewInfras(newReadWriter(t), newKVS(t))
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
				req := tc.arrange(t, ctx, infras, user.ID)
				server := newSessionServer(t, infras)

				// act
				client := authv1connect.NewSessionServiceClient(http.DefaultClient, server.URL)
				connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
				got, err := client.Refresh(ctx, connReq)

				// assert
				tc.assert(t, got, err)
			})
		}
	})

	t.Run("in cookie", func(t *testing.T) {
		tsc := []struct {
			name    string
			arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *http.Cookie
			assert  func(t *testing.T, got *connect.Response[authv1.RefreshResponse], err error)
		}{
			{
				name: "success",
				arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *http.Cookie {
					session := authModel.MustNewSession(userID)
					require.NoError(t, authRepository.NewSession(infras.KVS).Create(ctx, session))
					return &http.Cookie{
						Name:  "refresh_token",
						Value: session.Tokens.Refresh.Value,
					}
				},
				assert: func(t *testing.T, got *connect.Response[authv1.RefreshResponse], err error) {
					require.NoError(t, err)
					assert.NotEmpty(t, got.Msg.Tokens.Access.Value)
					assert.NotEmpty(t, got.Msg.Tokens.Access.ExpiresAt)
					assert.NotEmpty(t, got.Msg.Tokens.Refresh.Value)
					assert.NotEmpty(t, got.Msg.Tokens.Refresh.ExpiresAt)
				},
			},
			{
				name: "fail (refresh token not found)",
				arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *http.Cookie {
					return &http.Cookie{
						Name:  "refresh_token",
						Value: "",
					}
				},
				assert: func(t *testing.T, got *connect.Response[authv1.RefreshResponse], err error) {
					require.Error(t, err)
					assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
					connErr := new(connect.Error)
					require.ErrorAs(t, err, &connErr)
					require.Len(t, connErr.Details(), 1)
					detail := either.Must(connErr.Details()[0].Value())
					if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
						expectedMsg := i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseErrorInvalid_refresh_token})
						assert.Equal(t, expectedMsg, errMsg.Message)
					} else {
						require.Failf(t, "unexpected detail type", "got=%T", detail)
					}
				},
			},
			{
				name: "fail (invalid refresh token)",
				arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *http.Cookie {
					session := authModel.MustNewSession(userID)
					require.NoError(t, authRepository.NewSession(infras.KVS).Create(ctx, session))
					return &http.Cookie{
						Name:  "refresh_token",
						Value: session.Tokens.Refresh.Value + "invalid",
					}
				},
				assert: func(t *testing.T, got *connect.Response[authv1.RefreshResponse], err error) {
					require.Error(t, err)
					assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
					connErr := new(connect.Error)
					require.ErrorAs(t, err, &connErr)
					require.Len(t, connErr.Details(), 1)
					detail := either.Must(connErr.Details()[0].Value())
					if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
						expectedMsg := i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseErrorInvalid_refresh_token})
						assert.Equal(t, expectedMsg, errMsg.Message)
					} else {
						require.Failf(t, "unexpected detail type", "got=%T", detail)
					}
				},
			},
		}

		for _, tc := range tsc {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				// arrange
				ctx := context.Background()
				infras := di.NewInfras(newReadWriter(t), newKVS(t))
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
				cookie := tc.arrange(t, ctx, infras, user.ID)
				server := newSessionServer(t, infras)

				// act
				client := authv1connect.NewSessionServiceClient(http.DefaultClient, server.URL)
				req := &authv1.RefreshRequest{}
				connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
				connReq.Header().Add("Cookie", cookie.String())
				got, err := client.Refresh(ctx, connReq)

				// assert
				tc.assert(t, got, err)
			})
		}
	})
}

func TestSession_SignOut(t *testing.T) {
	t.Parallel()

	user := userFixture.User(nil)

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.SignOutRequest
		assert  func(t *testing.T, got *connect.Response[authv1.SignOutResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.SignOutRequest {
				session := authModel.MustNewSession(userID)
				require.NoError(t, authRepository.NewSession(infras.KVS).Create(ctx, session))
				return &authv1.SignOutRequest{
					AccessToken:  session.Tokens.Access.Value,
					RefreshToken: session.Tokens.Refresh.Value,
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.SignOutResponse], err error) {
				require.NoError(t, err)
				cookies := got.Header().Values("Set-Cookie")
				assert.Len(t, cookies, 2)
				assert.Contains(t, cookies, "access_token=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly; Secure")
				assert.Contains(t, cookies, "refresh_token=; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly; Secure")
			},
		},
		{
			name: "fail (session not found)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *authv1.SignOutRequest {
				return &authv1.SignOutRequest{
					AccessToken:  gofakeit.UUID(),
					RefreshToken: gofakeit.UUID(),
				}
			},
			assert: func(t *testing.T, got *connect.Response[authv1.SignOutResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
					expectedMsg := i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.AuthUsecaseErrorInvalid_access_refresh_token})
					assert.Equal(t, expectedMsg, errMsg.Message)
				} else {
					require.Failf(t, "unexpected detail type", "got=%T", detail)
				}
			},
		},
	}

	for _, tc := range tsc {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// arrange
			ctx := context.Background()
			infras := di.NewInfras(newReadWriter(t), newKVS(t))
			require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.ID)
			server := newSessionServer(t, infras)

			// act
			client := authv1connect.NewSessionServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.SignOut(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newSessionServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitAuthHandlers(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).Session
		return authv1connect.NewSessionServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
