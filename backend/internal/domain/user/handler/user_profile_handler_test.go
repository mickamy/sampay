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
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/ptr"
	"mickamy.com/sampay/test/connecttest"
)

func TestUserProfile_UpdateUserProfile(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.UpdateUserProfileRequest
		assert  func(t *testing.T, got *connect.Response[userv1.UpdateUserProfileResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.UpdateUserProfileRequest {
				return &userv1.UpdateUserProfileRequest{
					Name: "name",
					Bio:  ptr.Of("bio"),
					Image: &commonv1.S3Object{
						Bucket: gofakeit.GlobalFaker.Username(),
						Key:    gofakeit.GlobalFaker.UUID(),
					},
				}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.UpdateUserProfileResponse], err error) {
				require.NoError(t, err)
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
			})
			require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.ID)
			server := newUserProfileServer(t, infras)

			// act
			client := userv1connect.NewUserProfileServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.UpdateUserProfile(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestUserProfile_DeleteUserProfileImage(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.DeleteUserProfileImageRequest
		assert  func(t *testing.T, got *connect.Response[userv1.DeleteUserProfileImageResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras, userID string) *userv1.DeleteUserProfileImageRequest {
				return &userv1.DeleteUserProfileImageRequest{}
			},
			assert: func(t *testing.T, got *connect.Response[userv1.DeleteUserProfileImageResponse], err error) {
				require.NoError(t, err)
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
			})
			require.NoError(t, infras.Writer.DB.WithContext(ctx).Create(&user).Error)
			req := tc.arrange(t, ctx, infras, user.ID)
			server := newUserProfileServer(t, infras)

			// act
			client := userv1connect.NewUserProfileServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.DeleteUserProfileImage(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newUserProfileServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitUserHandler(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).UserProfile
		return userv1connect.NewUserProfileServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
