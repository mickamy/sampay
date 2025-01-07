package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/ptr"
)

func UserProfile(setter func(m *model.UserProfile)) model.UserProfile {
	m := model.UserProfile{
		Name: gofakeit.GlobalFaker.Name(),
		Bio:  ptr.Of(gofakeit.GlobalFaker.Sentence(20)),
	}

	if setter != nil {
		setter(&m)
	}

	return m
}
