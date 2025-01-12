package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"mickamy.com/sampay/internal/domain/user/model"
)

func UserLinkDisplayAttribute(setter func(m *model.UserLinkDisplayAttribute)) model.UserLinkDisplayAttribute {
	m := model.UserLinkDisplayAttribute{
		Name:         gofakeit.GlobalFaker.BeerName(),
		DisplayOrder: 0,
	}

	if setter != nil {
		setter(&m)
	}

	return m
}
