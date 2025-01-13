package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/user/v1/userv1connect"
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"
	"connectrpc.com/connect"
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

func TestUser_GetMe(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.GetMeRequest
		assert  func(t *testing.T, got *connect.Response[userv1.GetMeResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.GetMeRequest {
				return &userv1.GetMeRequest{}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.GetMeResponse], err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Msg.User)
				assert.NotEmpty(t, got.Msg.User.Profile)
				require.NotEmpty(t, got.Msg.User.Links)
				assert.NotEmpty(t, got.Msg.User.Links[0])
				assert.NotEmpty(t, got.Msg.User.Links[0].QrCodeUrl)
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
						m.QRCode = ptr.Of(commonFixture.S3Object(nil))
					}),
				}
			})
			require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.ID)
			server := newUserServer(t, infras)

			// act
			client := userv1connect.NewUserServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.GetMe(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestUser_GetUser(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, slug string) *userv1.GetUserRequest
		assert  func(t *testing.T, got *connect.Response[userv1.GetUserResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, slug string) *userv1.GetUserRequest {
				return &userv1.GetUserRequest{
					Slug: slug,
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.GetUserResponse], err error) {
				require.NoError(t, err)
				assert.NotEmpty(t, got.Msg.User)
				assert.NotEmpty(t, got.Msg.User.Profile)
				assert.NotEmpty(t, got.Msg.User.Profile.ImageUrl)
				require.NotEmpty(t, got.Msg.User.Links)
				assert.NotEmpty(t, got.Msg.User.Links[0])
				assert.NotEmpty(t, got.Msg.User.Links[0].QrCodeUrl)
			},
		},
		{
			name: "fail (empty user slug)",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, slug string) *userv1.GetUserRequest {
				return &userv1.GetUserRequest{
					Slug: "",
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.GetUserResponse], err error) {
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.ErrorMessage); ok {
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: "user.usecase.get_user.error.not_found"}), errMsg.Message)
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
			user := userFixture.User(func(m *userModel.User) {
				m.Profile = userFixture.UserProfile(func(m *userModel.UserProfile) {
					m.SetImage(ptr.Of(commonFixture.S3Object(nil)))
				})
				m.Links = []userModel.UserLink{
					userFixture.UserLink(func(m *userModel.UserLink) {
						m.QRCode = ptr.Of(commonFixture.S3Object(nil))
					}),
				}
			})
			require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.Slug)
			server := newUserServer(t, infras)

			// act
			client := userv1connect.NewUserServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.GetUser(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newUserServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitUserHandler(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).User
		return userv1connect.NewUserServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
