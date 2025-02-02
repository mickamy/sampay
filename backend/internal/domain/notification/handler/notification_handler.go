package handler

import (
	"context"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/notification/v1/notificationv1connect"
	notificationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/notification/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/domain/notification/model"
	"mickamy.com/sampay/internal/domain/notification/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/slices"
)

type Notification struct {
	list usecase.ListNotifications
}

func NewNotification(
	list usecase.ListNotifications,
) *Notification {
	return &Notification{
		list: list,
	}
}

func (h *Notification) ListNotifications(
	ctx context.Context,
	req *connect.Request[notificationv1.ListNotificationsRequest],
) (*connect.Response[notificationv1.ListNotificationsResponse], error) {
	out, err := h.list.Do(ctx, usecase.ListNotificationsInput{})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&notificationv1.ListNotificationsResponse{
		Notifications: slices.Map(out.Notifications, func(m model.Notification) *notificationv1.Notification {
			return &notificationv1.Notification{
				Id:        m.ID,
				Subject:   m.Subject,
				Body:      m.Body,
				CreatedAt: timestamppb.New(m.CreatedAt),
			}
		}),
	})
	return res, nil
}

func (h *Notification) ReadNotification(
	ctx context.Context,
	req *connect.Request[notificationv1.ReadNotificationRequest],
) (*connect.Response[notificationv1.ReadNotificationResponse], error) {
	panic("not implemented")
}

var _ notificationv1connect.NotificationServiceHandler = (*Notification)(nil)
