package fixture

import (
	"github.com/brianvoe/gofakeit/v7"

	"mickamy.com/sampay/internal/domain/common/model"
)

func S3Object(setter func(m *model.S3Object)) model.S3Object {
	m := model.S3Object{
		Bucket:      gofakeit.GlobalFaker.ProductName(),
		Key:         gofakeit.GlobalFaker.UUID(),
		ContentType: model.MustNewContentType(ContentType()),
	}

	if setter != nil {
		setter(&m)
	}

	return m
}

func ContentType() string {
	return gofakeit.GlobalFaker.RandomString([]string{
		"audio/mpeg",
		"image/bmp",
		"image/gif",
		"image/jpeg",
		"image/jpg",
		"image/png",
		"video/mp4",
		"video/mpeg",
	})
}
