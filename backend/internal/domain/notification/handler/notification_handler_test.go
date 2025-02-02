package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/notification/v1/notificationv1connect"
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
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
	"mickamy.com/sampay/internal/lib/either"
	"mickamy.com/sampay/internal/misc/i18n"
	"mickamy.com/sampay/test/connecttest"
)

func TestNotification_ListNotifications(t *testing.T) {
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
				m1 := fixture.Notification(func(m *model.Notification) {
					m.UserID = userID
				})
				m2 := fixture.Notification(func(m *model.Notification) {
					m.UserID = userID
					m.ReadStatus = fixture.NotificationReadStatusRead(func(m *model.NotificationReadStatus) {
						m.UserID = userID
					})
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&m1).Error)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&m2).Error)
				return &notificationv1.ListNotificationsRequest{
					Page: &commonv1.Page{
						Index: 0,
						Limit: 10,
					},
				}
			},
			assert: func(t *testing.T, got *connect.Response[notificationv1.ListNotificationsResponse], err error) {
				require.NoError(t, err)
				require.Len(t, got.Msg.Notifications, 2)
				assert.NotEmpty(t, got.Msg.Notifications[0].Id)
				assert.NotEmpty(t, got.Msg.Notifications[0].Subject)
				assert.NotEmpty(t, got.Msg.Notifications[0].Body)
				assert.NotEmpty(t, got.Msg.Notifications[0].CreatedAt)
				assert.NotEmpty(t, got.Msg.Notifications[0].ReadAt)
			},
		},
		{
			name: "fail (page not set)",
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
				require.Error(t, err)
				assert.Equalf(t, connect.CodeInvalidArgument, connect.CodeOf(err), "code=%s", connect.CodeOf(err).String())
				connErr := new(connect.Error)
				require.ErrorAs(t, err, &connErr)
				require.Len(t, connErr.Details(), 1)
				detail := either.Must(connErr.Details()[0].Value())
				if errMsg, ok := detail.(*commonv1.BadRequestError); ok {
					require.Len(t, errMsg.FieldViolations, 1)
					require.Equal(t, "page", errMsg.FieldViolations[0].Field)
					require.Len(t, errMsg.FieldViolations[0].Descriptions, 1)
					require.Equal(t, i18n.MustJapaneseMessage(i18n.Config{MessageID: i18n.CommonHandlerErrorInvalid_page}), errMsg.FieldViolations[0].Descriptions[0])
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

func TestNotification_ReadNotification(t *testing.T) {
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

func TestNotification_UnreadNotificationsCount(t *testing.T) {
	t.Parallel()

	tsc := []struct {
		name    string
		arrange func(t *testing.T, ctx context.Context, infras di.Infras) *notificationv1.CountUnreadNotificationRequest
		assert  func(t *testing.T, got *connect.Response[notificationv1.CountUnreadNotificationResponse], err error)
	}{
		{
			name: "success",
			arrange: func(t *testing.T, ctx context.Context, infras di.Infras) *notificationv1.CountUnreadNotificationRequest {
				userID := contexts.MustAuthenticatedUserID(ctx)
				m1 := fixture.Notification(func(m *model.Notification) {
					m.UserID = userID
				})
				m2 := fixture.Notification(func(m *model.Notification) {
					m.UserID = userID
					m.ReadStatus = fixture.NotificationReadStatusRead(func(m *model.NotificationReadStatus) {
						m.UserID = userID
					})
				})
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&m1).Error)
				require.NoError(t, infras.Writer.WithContext(ctx).Create(&m2).Error)
				return &notificationv1.CountUnreadNotificationRequest{}
			},
			assert: func(t *testing.T, got *connect.Response[notificationv1.CountUnreadNotificationResponse], err error) {
				require.NoError(t, err)
				assert.Equal(t, int32(1), got.Msg.Count)
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
			got, err := client.CountUnreadNotification(ctx, connReq)

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
