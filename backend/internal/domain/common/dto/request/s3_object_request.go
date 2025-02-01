package request

import (
	commonv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"

	commonModel "mickamy.com/sampay/internal/domain/common/model"
)

func NewS3Object(pb *commonv1.S3Object) (*commonModel.S3Object, error) {
	if pb == nil {
		return nil, nil
	}

	contentType, err := commonModel.NewContentType(pb.ContentType)
	if err != nil {
		return nil, err
	}

	return &commonModel.S3Object{
		Bucket:      pb.Bucket,
		Key:         pb.Key,
		ContentType: contentType,
	}, nil
}
