package usecase

import (
	"context"
	"strings"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/storage/model"
	"github.com/mickamy/sampay/internal/domain/storage/repository"
	s3client "github.com/mickamy/sampay/internal/infra/aws/s3"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/lib/ulid"
)

type GetUploadURLInput struct {
	Path string
}

type GetUploadURLOutput struct {
	UploadURL  string
	S3ObjectID string
}

type GetUploadURL interface {
	Do(ctx context.Context, input GetUploadURLInput) (GetUploadURLOutput, error)
}

type getUploadURL struct {
	_         GetUploadURL        `inject:"returns"`
	_         *di.Infra           `inject:"param"`
	s3        s3client.Client     `inject:""`
	writer    *database.Writer    `inject:""`
	s3ObjRepo repository.S3Object `inject:""`
}

func (uc *getUploadURL) Do(ctx context.Context, input GetUploadURLInput) (GetUploadURLOutput, error) {
	if err := validatePath(input.Path); err != nil {
		return GetUploadURLOutput{}, err
	}

	bucket := config.AWS().S3PublicBucket
	key := input.Path

	uploadURL, err := uc.s3.PresignPutObject(ctx, bucket, key)
	if err != nil {
		return GetUploadURLOutput{}, errx.
			Wrap(err, "failed to generate presigned URL", "bucket", bucket, "key", key).
			WithCode(errx.Internal)
	}

	obj := model.S3Object{
		ID:     ulid.New(),
		Bucket: bucket,
		Key:    key,
	}

	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		if err := uc.s3ObjRepo.WithTx(tx).Create(ctx, &obj); err != nil {
			return errx.
				Wrap(err, "failed to create s3 object record").
				WithCode(errx.Internal)
		}
		return nil
	}); err != nil {
		//nolint:wrapcheck // errors from transaction callback are already wrapped inside
		return GetUploadURLOutput{}, err
	}

	return GetUploadURLOutput{
		UploadURL:  uploadURL,
		S3ObjectID: obj.ID,
	}, nil
}

func validatePath(path string) error {
	if path == "" {
		return errx.New("path is required").WithCode(errx.InvalidArgument)
	}
	if strings.Contains(path, "..") {
		return errx.New("path must not contain '..'").WithCode(errx.InvalidArgument)
	}
	if strings.HasPrefix(path, "/") {
		return errx.New("path must not start with '/'").WithCode(errx.InvalidArgument)
	}
	return nil
}
