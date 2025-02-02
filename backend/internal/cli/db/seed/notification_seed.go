package seed

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/notification/fixture"
	"mickamy.com/sampay/internal/domain/notification/model"
	userModel "mickamy.com/sampay/internal/domain/user/model"
)

func seedNotification(ctx context.Context, writer *database.Writer, env config.Env) error {
	if env != config.Development {
		return nil
	}

	var me userModel.User
	if err := writer.WithContext(ctx).First(&me, "slug = ?", "mickamy").Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("do not seed notification because user not found")
			return nil
		}
		return fmt.Errorf("failed to find user: %w", err)
	}

	var count int64
	if err := writer.WithContext(ctx).Model(&model.Notification{}).Where("user_id = ?", me.ID).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count notifications: %w", err)
	}

	if count >= 10 {
		fmt.Println("do not seed notification because user already has 10 notifications")
		return nil
	}

	var notifications []model.Notification
	for i := 0; i < 10; i++ {
		notifications = append(notifications, fixture.NotificationMessageReceived(func(m *model.Notification) {
			m.UserID = me.ID
		}))
	}
	if err := writer.WithContext(ctx).Create(&notifications).Error; err != nil {
		return fmt.Errorf("failed to create notifications: %w", err)
	}

	return nil
}
