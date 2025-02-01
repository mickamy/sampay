package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"mickamy.com/sampay/internal/domain/message/model"
)

func Message(setter func(m *model.Message)) model.Message {
	m := model.Message{
		SenderName: gofakeit.GlobalFaker.Name(),
		Content:    gofakeit.GlobalFaker.Sentence(10),
	}

	if setter != nil {
		setter(&m)
	}

	return m
}
