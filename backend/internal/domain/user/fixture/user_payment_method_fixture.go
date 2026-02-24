package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/mickamy/sampay/internal/domain/user/model"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

func UserPaymentMethod(setter func(m *model.UserPaymentMethod)) model.UserPaymentMethod {
	m := model.UserPaymentMethod{
		ID:           ulid.New(),
		UserID:       ulid.New(),
		Type:         "paypay",
		URL:          gofakeit.URL(),
		DisplayOrder: 0,
	}
	if setter != nil {
		setter(&m)
	}
	return m
}
