package fixture

import (
	"github.com/mickamy/sampay/internal/domain/storage/model"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

func S3Object(setter func(m *model.S3Object)) model.S3Object {
	m := model.S3Object{
		ID:     ulid.New(),
		Bucket: "test-bucket",
		Key:    "test-key/" + ulid.New() + ".png",
	}
	if setter != nil {
		setter(&m)
	}
	return m
}
