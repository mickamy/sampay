package fixture

import (
	"time"

	"mickamy.com/sampay/internal/domain/notification/model"
	"mickamy.com/sampay/internal/lib/ptr"
)

func NotificationReadStatus(setter func(m *model.NotificationReadStatus)) model.NotificationReadStatus {
	m := model.NotificationReadStatus{}

	if setter != nil {
		setter(&m)
	}

	return m
}

func NotificationReadStatusRead(setter func(m *model.NotificationReadStatus)) model.NotificationReadStatus {
	m := NotificationReadStatus(func(m *model.NotificationReadStatus) {
		m.ReadAt = ptr.Of(time.Now())
	})

	if setter != nil {
		setter(&m)
	}

	return m
}
