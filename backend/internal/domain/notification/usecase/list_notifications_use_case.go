package usecase

import (
	"context"
	"fmt"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/notification/model"
	notificationRepository "mickamy.com/sampay/internal/domain/notification/repository"
	"mickamy.com/sampay/internal/lib/contexts"
)

type ListNotificationsInput struct {
}

type ListNotificationsOutput struct {
	Notifications []model.Notification
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
	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		var err error
		notifications, err = uc.notificationRepo.WithTx(tx.ReaderDB()).ListByUserID(ctx, contexts.MustAuthenticatedUserID(ctx))
		if err != nil {
			return fmt.Errorf("failed to list notifications: %w", err)
		}

		return nil
	}); err != nil {
		return ListNotificationsOutput{}, err
	}

	return ListNotificationsOutput{Notifications: notifications}, nil
}
