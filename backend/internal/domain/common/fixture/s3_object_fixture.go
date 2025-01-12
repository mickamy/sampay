package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"mickamy.com/sampay/internal/domain/common/model"
)

func S3Object(setter func(m *model.S3Object)) model.S3Object {
	m := model.S3Object{
		Bucket: gofakeit.GlobalFaker.ProductName(),
		Key:    gofakeit.GlobalFaker.UUID(),
	}

	if setter != nil {
		setter(&m)
	}

	return m
}
