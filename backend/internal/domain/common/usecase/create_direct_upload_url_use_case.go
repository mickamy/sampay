package usecase

import (
	"context"

	commonModel "mickamy.com/sampay/internal/domain/common/model"
	"mickamy.com/sampay/internal/lib/aws/s3"
)

type CreateDirectUploadURLInput struct {
	commonModel.S3Object
}

type CreateDirectUploadURLOutput struct {
	URL string
}

//go:generate mockgen -source=$GOFILE -destination=./mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
type CreateDirectUploadURL interface {
	Do(ctx context.Context, input CreateDirectUploadURLInput) (CreateDirectUploadURLOutput, error)
}

type createDirectUploadURL struct {
	s3 s3.Client
}

func NewCreateDirectUploadURL(
	s3 s3.Client,
) CreateDirectUploadURL {
	return &createDirectUploadURL{
		s3: s3,
	}
}

func (uc *createDirectUploadURL) Do(ctx context.Context, input CreateDirectUploadURLInput) (CreateDirectUploadURLOutput, error) {
	url, err := uc.s3.GeneratePresignedURL(ctx, input.Bucket, input.Key, 60)
	if err != nil {
		return CreateDirectUploadURLOutput{}, err
	}

	return CreateDirectUploadURLOutput{URL: url}, nil
}
