package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"mickamy.com/sampay/internal/domain/user/model"
)

func User(setter func(m *model.User)) model.User {
	m := model.User{
		Slug: gofakeit.GlobalFaker.Username(),
	}

	if setter != nil {
		setter(&m)
	}

	return m
}
