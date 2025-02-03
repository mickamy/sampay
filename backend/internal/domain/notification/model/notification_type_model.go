package model

type NotificationType string

func (m NotificationType) String() string {
	return string(m)
}

const (
	NotificationTypeAnnouncement NotificationType = "announcement"
	NotificationTypeMessage      NotificationType = "message"
)
