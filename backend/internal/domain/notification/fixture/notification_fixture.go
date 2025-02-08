package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/mattn/go-gimei"

	commonFixture "mickamy.com/sampay/internal/domain/common/fixture"
	"mickamy.com/sampay/internal/domain/notification/model"
	"mickamy.com/sampay/internal/misc/i18n"
)

func Notification(setter func(m *model.Notification)) model.Notification {
	m := model.Notification{
		Type:    NotificationType(),
		Subject: gofakeit.GlobalFaker.Sentence(2),
		Body:    gofakeit.GlobalFaker.Sentence(20),
	}

	if setter != nil {
		setter(&m)
	}

	return m
}

func NotificationMessageReceived(setter func(m *model.Notification)) model.Notification {
	m := Notification(func(m *model.Notification) {
		name := gimei.NewName()
		m.Subject = i18n.MustJapaneseMessage(
			i18n.Config{
				MessageID:    i18n.MessageUsecaseCreate_messageEmailSubject,
				TemplateData: map[string]string{"SenderName": name.Kanji()},
			},
		)
		m.Body = i18n.MustJapaneseMessage(
			i18n.Config{
				MessageID:    i18n.MessageUsecaseCreate_messageEmailBody,
				TemplateData: map[string]string{"SenderName": name.Kanji(), "Content": gofakeit.Sentence(2)},
			},
		)
	})

	if setter != nil {
		setter(&m)
	}

	return m
}

func NotificationType() model.NotificationType {
	return commonFixture.RandomStringer(model.NotificationTypeValues())
}
