package fixture

import (
	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

func User(setter func(m *model.User)) model.User {
	m := model.User{
		ID: ulid.New(),
	}
	if setter != nil {
		setter(&m)
	}
	return m
}
