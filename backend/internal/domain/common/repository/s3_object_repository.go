package repository

import (
	"context"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/common/model"
)

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type S3Object interface {
	Upsert(ctx context.Context, m *model.S3Object) error
	WithTx(tx *database.DB) S3Object
}

type s3Object struct {
	db *database.DB
}

func NewS3Object(db *database.DB) S3Object {
	return &s3Object{db: db}
}

func (repo *s3Object) Upsert(ctx context.Context, m *model.S3Object) error {
	return repo.db.WithContext(ctx).
		FirstOrCreate(m, "bucket = ? AND key = ?", m.Bucket, m.Key).
		Error
}

func (repo *s3Object) WithTx(tx *database.DB) S3Object {
	return &s3Object{db: tx}
}
