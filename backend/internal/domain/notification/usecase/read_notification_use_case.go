package usecase

import (
	"context"
	"fmt"

	notificationRepository "mickamy.com/sampay/internal/domain/notification/repository"
	"mickamy.com/sampay/internal/infra/storage/database"
)

type ReadNotificationInput struct {
	ID string
}

type ReadNotificationOutput struct {
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type ReadNotification interface {
	Do(ctx context.Context, input ReadNotificationInput) (ReadNotificationOutput, error)
}

type readNotification struct {
	writer           *database.Writer
	notificationRepo notificationRepository.Notification
}

func NewReadNotification(
	writer *database.Writer,
	notificationRepo notificationRepository.Notification,
) ReadNotification {
	return &readNotification{
		writer:           writer,
		notificationRepo: notificationRepo,
	}
}

func (uc *readNotification) Do(ctx context.Context, input ReadNotificationInput) (ReadNotificationOutput, error) {
	if err := uc.writer.WriterTransaction(ctx, func(tx database.WriterTransactional) error {
		notification, err := uc.notificationRepo.WithTx(tx.WriterDB()).Find(ctx, input.ID)
		if err != nil {
			return fmt.Errorf("failed to find notification: %w", err)
		}
		if notification == nil {
			return fmt.Errorf("notification not found: id=[%s]", input.ID)
		}

		notification.Read()
		if err := uc.notificationRepo.WithTx(tx.WriterDB()).Update(ctx, notification); err != nil {
			return fmt.Errorf("failed to update notification: %w", err)
		}

		return nil
	}); err != nil {
		return ReadNotificationOutput{}, err
	}

	return ReadNotificationOutput{}, nil
}
