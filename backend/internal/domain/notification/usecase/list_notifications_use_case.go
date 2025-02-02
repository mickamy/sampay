package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/notification/model"
	notificationRepository "mickamy.com/sampay/internal/domain/notification/repository"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/paging"
)

type ListNotificationsInput struct {
	Page paging.Page
}

type ListNotificationsOutput struct {
	Notifications []model.Notification
	NextPage      paging.NextPage
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type ListNotifications interface {
	Do(ctx context.Context, input ListNotificationsInput) (ListNotificationsOutput, error)
}

type listNotifications struct {
	reader           *database.Reader
	notificationRepo notificationRepository.Notification
}

func NewListNotifications(
	reader *database.Reader,
	notificationRepo notificationRepository.Notification,
) ListNotifications {
	return &listNotifications{
		reader:           reader,
		notificationRepo: notificationRepo,
	}
}

func (uc *listNotifications) Do(ctx context.Context, input ListNotificationsInput) (ListNotificationsOutput, error) {
	var notifications []model.Notification
	var nextPage paging.NextPage
	userID := contexts.MustAuthenticatedUserID(ctx)
	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		var err error
		notifications, err = uc.notificationRepo.WithTx(tx.ReaderDB()).ListByUserID(ctx, userID, notificationRepository.NotificationJoinReadStatus, notificationRepository.NotificationOrderByIDDesc, input.Page.Scope())
		if err != nil {
			return fmt.Errorf("failed to list notifications: %w", err)
		}

		count, err := uc.notificationRepo.WithTx(tx.ReaderDB()).CountByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to count notifications: %w", err)
		}

		nextPage = input.Page.Next(count)

		return nil
	}); err != nil {
		return ListNotificationsOutput{}, err
	}

	return ListNotificationsOutput{Notifications: notifications, NextPage: nextPage}, nil
}
