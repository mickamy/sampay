package model

import (
	"time"

	smodel "github.com/mickamy/sampay/internal/domain/storage/model"
)

//go:generate go tool ormgen -source=$GOFILE -destination=../query
type UserPaymentMethod struct {
	ID                string
	UserID            string
	Type              string
	URL               string
	QRCodeS3ObjectID  *string `db:",nullable"`
	DisplayOrder      int
	CreatedAt         time.Time
	UpdatedAt         time.Time

	QRCodeS3Object *smodel.S3Object `rel:"belongs_to,foreign_key:qr_code_s3_object_id"`
}
