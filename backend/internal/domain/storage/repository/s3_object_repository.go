package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/mickamy/ormgen/orm"

	"github.com/mickamy/sampay/internal/domain/storage/model"
	"github.com/mickamy/sampay/internal/domain/storage/query"
	"github.com/mickamy/sampay/internal/infra/storage/database"
)

type S3Object interface {
	Create(ctx context.Context, m *model.S3Object) error
	Get(ctx context.Context, id string) (model.S3Object, error)
	WithTx(tx *database.DB) S3Object
}

type s3Object struct {
	db *database.DB
}

func NewS3Object(db *database.DB) S3Object {
	return &s3Object{db: db}
}

func (repo *s3Object) Create(ctx context.Context, m *model.S3Object) error {
	if err := query.S3Objects(repo.db).Create(ctx, m); err != nil {
		return fmt.Errorf("repository: %w", err)
	}
	return nil
}

func (repo *s3Object) Get(ctx context.Context, id string) (model.S3Object, error) {
	m, err := query.S3Objects(repo.db).Where("id = ?", id).First(ctx)
	if errors.Is(err, orm.ErrNotFound) {
		return model.S3Object{}, database.ErrNotFound
	}
	if err != nil {
		return m, fmt.Errorf("repository: %w", err)
	}
	return m, nil
}

func (repo *s3Object) WithTx(tx *database.DB) S3Object {
	return &s3Object{db: tx}
}
