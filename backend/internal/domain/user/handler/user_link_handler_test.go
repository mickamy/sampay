package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/user/v1/userv1connect"
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/lib/ptr"
	"mickamy.com/sampay/internal/misc/i18n"
	"mickamy.com/sampay/test/connecttest"
)

func TestUserLink_CreateUserLink(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.CreateUserLinkRequest
		assert  func(t *testing.T, got *connect.Response[userv1.CreateUserLinkResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.CreateUserLinkRequest {
				m := userFixture.UserLink(func(m *userModel.UserLink) {
					m.DisplayAttribute = userFixture.UserLinkDisplayAttribute(nil)
				})
				return &userv1.CreateUserLinkRequest{
					ProviderType: m.ProviderType.String(),
					Uri:          m.URI,
					Name:         m.DisplayAttribute.Name,
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.CreateUserLinkResponse], err error) {
				require.NoError(t, err)
				assert.Empty(t, got.Msg.String())
			},
		},
		{
			name: "fail (empty provider type)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.CreateUserLinkRequest {
				m := userFixture.UserLink(func(m *userModel.UserLink) {
					m.DisplayAttribute = userFixture.UserLinkDisplayAttribute(nil)
				})
				return &userv1.CreateUserLinkRequest{
					ProviderType: "",
					Uri:          m.URI,
					Name:         m.DisplayAttribute.Name,
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.CreateUserLinkResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.BadRequestError); ok {
					require.Len(t, errMsg.FieldViolations, 1)
					require.Equal(t, "provider_type", errMsg.FieldViolations[0].Field)
					require.Len(t, errMsg.FieldViolations[0].Descriptions, 1)
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.UserHandlerUser_linkErrorInvalid_provider_type}), errMsg.FieldViolations[0].Descriptions[0])
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
			require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.ID)
			server := newUserLinkServer(t, infras)

			// act
			client := userv1connect.NewUserLinkServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.CreateUserLink(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestUserLink_ListUserLink(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.ListUserLinkRequest
		assert  func(t *testing.T, got *connect.Response[userv1.ListUserLinkResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.ListUserLinkRequest {
				m := userFixture.UserLink(func(m *userModel.UserLink) {
					m.UserID = userID
					m.DisplayAttribute = userFixture.UserLinkDisplayAttribute(nil)
				})
				require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&m).Error)
				return &userv1.ListUserLinkRequest{
					UserId: userID,
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.ListUserLinkResponse], err error) {
				require.NoError(t, err)
				require.NotEmpty(t, got.Msg.Links)
				require.Len(t, got.Msg.Links, 1)
				assert.NotEmpty(t, got.Msg.Links[0])
			},
		},
		{
			name: "empty user id",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.ListUserLinkRequest {
				return &userv1.ListUserLinkRequest{
					UserId: "",
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.ListUserLinkResponse], err error) {
				require.NoError(t, err)
				assert.Empty(t, got.Msg.Links)
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
			require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.ID)
			server := newUserLinkServer(t, infras)

			// act
			client := userv1connect.NewUserLinkServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.ListUserLink(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestUserLink_UpdateUserLink(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.UpdateUserLinkRequest
		assert  func(t *testing.T, got *connect.Response[userv1.UpdateUserLinkResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.UpdateUserLinkRequest {
				m := userFixture.UserLink(func(m *userModel.UserLink) {
					m.UserID = userID
					m.DisplayAttribute = userFixture.UserLinkDisplayAttribute(nil)
				})
				require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&m).Error)
				return &userv1.UpdateUserLinkRequest{
					Id:           m.ID,
					ProviderType: ptr.Of(m.ProviderType.String()),
					Uri:          &m.URI,
					Name:         ptr.Of("updated"),
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.UpdateUserLinkResponse], err error) {
				require.NoError(t, err)
				assert.Empty(t, got.Msg.String())
			},
		},
		{
			name: "fail (empty provider type)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.UpdateUserLinkRequest {
				return &userv1.UpdateUserLinkRequest{
					ProviderType: ptr.Of(""),
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.UpdateUserLinkResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.BadRequestError); ok {
					require.Len(t, errMsg.FieldViolations, 1)
					require.Equal(t, "provider_type", errMsg.FieldViolations[0].Field)
					require.Len(t, errMsg.FieldViolations[0].Descriptions, 1)
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.UserHandlerUser_linkErrorInvalid_provider_type}), errMsg.FieldViolations[0].Descriptions[0])
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
			require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.ID)
			server := newUserLinkServer(t, infras)

			// act
			client := userv1connect.NewUserLinkServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.UpdateUserLink(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestUserLink_UpdateUserLinkQRCode(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, linkID string) *userv1.UpdateUserLinkQRCodeRequest
		assert  func(t *testing.T, got *connect.Response[userv1.UpdateUserLinkQRCodeResponse], err error)
	}{
		{
			name: "success (image is nil)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, linkID string) *userv1.UpdateUserLinkQRCodeRequest {
				return &userv1.UpdateUserLinkQRCodeRequest{
					Id: linkID,
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.UpdateUserLinkQRCodeResponse], err error) {
				require.NoError(t, err)
				assert.Empty(t, got.Msg.String())
			},
		},
		{
			name: "success (image is not nil)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, linkID string) *userv1.UpdateUserLinkQRCodeRequest {
				return &userv1.UpdateUserLinkQRCodeRequest{
					Id: linkID,
					QrCode: &commonv1.S3Object{
						Bucket: gofakeit.GlobalFaker.ProductName(),
						Key:    gofakeit.GlobalFaker.UUID(),
					},
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.UpdateUserLinkQRCodeResponse], err error) {
				require.NoError(t, err)
				assert.Empty(t, got.Msg.String())
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
			user := userFixture.User(func(m *userModel.User) {
				m.Profile = userFixture.UserProfile(func(m *userModel.UserProfile) {
					m.SetImage(ptr.Of(commonFixture.S3Object(nil)))
				})
				m.Links = []userModel.UserLink{
					userFixture.UserLink(func(m *userModel.UserLink) {
						m.SetQRCode(ptr.Of(commonFixture.S3Object(nil)))
					}),
				}
			})
			require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.Links[0].ID)
			server := newUserLinkServer(t, infras)

			// act
			client := userv1connect.NewUserLinkServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.UpdateUserLinkQRCode(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestUserLink_DeleteUserLink(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.DeleteUserLinkRequest
		assert  func(t *testing.T, got *connect.Response[userv1.DeleteUserLinkResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.DeleteUserLinkRequest {
				m := userFixture.UserLink(func(m *userModel.UserLink) {
					m.UserID = userID
					m.DisplayAttribute = userFixture.UserLinkDisplayAttribute(nil)
				})
				require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&m).Error)
				return &userv1.DeleteUserLinkRequest{
					Id: m.ID,
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.DeleteUserLinkResponse], err error) {
				require.NoError(t, err)
				assert.Empty(t, got.Msg.String())
			},
		},
		{
			name: "empty id",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.DeleteUserLinkRequest {
				return &userv1.DeleteUserLinkRequest{
					Id: "",
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.DeleteUserLinkResponse], err error) {
				require.NoError(t, err)
				assert.Empty(t, got.Msg.String())
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
			require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.ID)
			server := newUserLinkServer(t, infras)

			// act
			client := userv1connect.NewUserLinkServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.DeleteUserLink(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newUserLinkServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitUserHandler(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).UserLink
		return userv1connect.NewUserLinkServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
