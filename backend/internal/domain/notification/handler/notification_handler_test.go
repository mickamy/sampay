package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/notification/v1/notificationv1connect"
	notificationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/notification/v1"
	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mickamy.com/sampay/internal/di"
	authModel "mickamy.com/sampay/internal/domain/auth/model"
	"mickamy.com/sampay/internal/domain/notification/fixture"
	"mickamy.com/sampay/internal/domain/notification/model"
	userFixture "mickamy.com/sampay/internal/domain/user/fixture"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/test/connecttest"
)

func TestMessage_SendMessage(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras) *notificationv1.ListNotificationsRequest
		assert  func(t *testing.T, got *connect.Response[notificationv1.ListNotificationsResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *notificationv1.ListNotificationsRequest {
				userID := contexts.MustAuthenticatedUserID(ctx)
				ms := []model.Notification{fixture.Notification(func(m *model.Notification) {
					m.UserID = userID
				}), fixture.Notification(func(m *model.Notification) {
					m.UserID = userID
					m.ReadStatus = fixture.NotificationReadStatusRead(func(m *model.NotificationReadStatus) {
						m.UserID = userID
					})
				})}
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&ms).Error)
				return &notificationv1.ListNotificationsRequest{}
			},
			assert: func(t *testing.T, got *connect.Response[notificationv1.ListNotificationsResponse], err error) {
				require.NoError(t, err)
				assert.Len(t, got.Msg.Notifications, 2)
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
			server := newMessageServer(t, infras)
			user := userFixture.User(nil)
			require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
			ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
			req := tc.arrange(t, ctx, infras)

			// act
			client := notificationv1connect.NewNotificationServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.ListNotifications(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func TestMessage_ReadMessage(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras) *notificationv1.ReadNotificationRequest
		assert  func(t *testing.T, got *connect.Response[notificationv1.ReadNotificationResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *notificationv1.ReadNotificationRequest {
				userID := contexts.MustAuthenticatedUserID(ctx)
				m := fixture.Notification(func(m *model.Notification) {
					m.UserID = userID
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&m).Error)
				return &notificationv1.ReadNotificationRequest{
					Id: m.ID,
				}
			},
			assert: func(t *testing.T, got *connect.Response[notificationv1.ReadNotificationResponse], err error) {
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
			server := newMessageServer(t, infras)
			user := userFixture.User(nil)
			require.NoError(t, infras.Writer.WithContext(ctx).Create(&user).Error)
			ctx = contexts.SetAuthenticatedUserID(ctx, user.ID)
			req := tc.arrange(t, ctx, infras)

			// act
			client := notificationv1connect.NewNotificationServiceClient(http.DefaultClient, server.URL)
			connReq := connecttest.NewAuthenticatedRequest(t, ctx, req, nil, authModel.MustNewSession(user.ID), infras.KVS)
			got, err := client.ReadNotification(ctx, connReq)

			// assert
			tc.assert(t, got, err)
		})
	}
}

func newMessageServer(t *testing.T, infras di.Infras) *httptest.Server {
	return connecttest.NewServer(t, infras, func(interceptors []connect.Interceptor) (string, http.Handler) {
		h := di.InitNotificationHandlers(infras.Writer.DB, infras.ReadWriter, infras.Writer, infras.Reader, infras.KVS).Notification
		return notificationv1connect.NewNotificationServiceHandler(h, connect.WithInterceptors(interceptors...))
	})
}
