package usecase

import (
	"context"
	"fmt"

	notificationRepository "mickamy.com/sampay/internal/domain/notification/repository"
	"mickamy.com/sampay/internal/infra/storage/database"
	"mickamy.com/sampay/internal/lib/contexts"
)

type CountUnreadNotificationsInput struct {
}

type CountUnreadNotificationsOutput struct {
	Count int
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CountUnreadNotifications interface {
	Do(ctx context.Context, input CountUnreadNotificationsInput) (CountUnreadNotificationsOutput, error)
}

type countUnreadNotifications struct {
	reader           *database.Reader
	notificationRepo notificationRepository.Notification
}

func NewCountUnreadNotifications(
	reader *database.Reader,
	notificationRepo notificationRepository.Notification,
) CountUnreadNotifications {
	return &countUnreadNotifications{
		reader:           reader,
		notificationRepo: notificationRepo,
	}
}

func (uc *countUnreadNotifications) Do(ctx context.Context, input CountUnreadNotificationsInput) (CountUnreadNotificationsOutput, error) {
	var count int
	userID := contexts.MustAuthenticatedUserID(ctx)
	if err := uc.reader.ReaderTransaction(ctx, func(tx database.ReaderTransactional) error {
		var err error
		count, err = uc.notificationRepo.WithTx(tx.ReaderDB()).CountByUserID(ctx, userID, notificationRepository.NotificationUnread)
		if err != nil {
			return fmt.Errorf("failed to count notifications: %w", err)
		}

		return nil
	}); err != nil {
		return CountUnreadNotificationsOutput{}, err
	}

	return CountUnreadNotificationsOutput{Count: count}, nil
}
