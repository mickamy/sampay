package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

func EndUser(setter func(m *model.EndUser)) model.EndUser {
	m := model.EndUser{
		UserID: ulid.New(),
		Slug:   gofakeit.Username(),
	}
	if setter != nil {
		setter(&m)
	}
	return m
}
