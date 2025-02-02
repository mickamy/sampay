package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"mickamy.com/sampay/internal/domain/notification/model"
)

func Notification(setter func(m *model.Notification)) model.Notification {
	m := model.Notification{
		Subject: gofakeit.GlobalFaker.Sentence(2),
		Body:    gofakeit.GlobalFaker.Sentence(20),
	}

	if setter != nil {
		setter(&m)
	}

	return m
}
