package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/notification/v1/notificationv1connect"
	notificationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/notification/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"
	"google.golang.org/protobuf/types/known/timestamppb"

	commonRequest "mickamy.com/sampay/internal/domain/common/dto/request"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/domain/notification/model"
	"mickamy.com/sampay/internal/domain/notification/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/operator"
	"mickamy.com/sampay/internal/lib/slices"
	"mickamy.com/sampay/internal/misc/i18n"
)

type Notification struct {
	list  usecase.ListNotifications
	read  usecase.ReadNotification
	count usecase.CountUnreadNotifications
}

func NewNotification(
	list usecase.ListNotifications,
	read usecase.ReadNotification,
	count usecase.CountUnreadNotifications,
) *Notification {
	return &Notification{
		list:  list,
		read:  read,
		count: count,
	}
}

func (h *Notification) ListNotifications(
	ctx context.Context,
	req *connect.Request[notificationv1.ListNotificationsRequest],
) (*connect.Response[notificationv1.ListNotificationsResponse], error) {
	lang := contexts.MustLanguage(ctx)

	page := commonRequest.NewPage(req.Msg.Page)
	if page == nil {
		return nil, commonResponse.NewBadRequest(errors.New("invalid page")).
			WithFieldViolation("page", i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.CommonHandlerErrorInvalid_page})).
			AsConnectError()
	}

	out, err := h.list.Do(ctx, usecase.ListNotificationsInput{
		Page: *page,
	})
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
				ReadAt: operator.TernaryFunc(
					m.ReadStatus.ReadAt != nil,
					func() *timestamppb.Timestamp {
						return timestamppb.New(*m.ReadStatus.ReadAt)
					}, func() *timestamppb.Timestamp {
						return nil
					},
				),
			}
		}),
	})
	return res, nil
}

func (h *Notification) ReadNotification(
	ctx context.Context,
	req *connect.Request[notificationv1.ReadNotificationRequest],
) (*connect.Response[notificationv1.ReadNotificationResponse], error) {
	_, err := h.read.Do(ctx, usecase.ReadNotificationInput{
		ID: req.Msg.Id,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&notificationv1.ReadNotificationResponse{})
	return res, nil
}

func (h *Notification) UnreadNotificationsCount(ctx context.Context, req *connect.Request[notificationv1.UnreadNotificationsCountRequest]) (*connect.Response[notificationv1.UnreadNotificationsCountResponse], error) {
	out, err := h.count.Do(ctx, usecase.CountUnreadNotificationsInput{})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&notificationv1.UnreadNotificationsCountResponse{
		Count: int32(out.Count),
	})
	return res, nil
}

var _ notificationv1connect.NotificationServiceHandler = (*Notification)(nil)
