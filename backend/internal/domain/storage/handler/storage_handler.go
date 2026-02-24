package handler

import (
	"context"

	"connectrpc.com/connect"

	v1 "github.com/mickamy/sampay/gen/storage/v1"
	"github.com/mickamy/sampay/gen/storage/v1/storagev1connect"
	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/storage/usecase"
	"github.com/mickamy/sampay/internal/lib/logger"
)

var _ storagev1connect.StorageServiceHandler = (*Storage)(nil)

type Storage struct {
	_            *di.Infra            `inject:"param"`
	getUploadURL usecase.GetUploadURL `inject:""`
}

func (h *Storage) GetUploadURL(
	ctx context.Context, r *connect.Request[v1.GetUploadURLRequest],
) (*connect.Response[v1.GetUploadURLResponse], error) {
	out, err := h.getUploadURL.Do(ctx, usecase.GetUploadURLInput{
		Path: r.Msg.GetPath(),
	})
	if err != nil {
		logger.Error(ctx, "failed to execute use-case", "err", err)
		return nil, err //nolint:wrapcheck // use-case errors are already wrapped with errx
	}

	return connect.NewResponse(&v1.GetUploadURLResponse{
		UploadUrl:  out.UploadURL,
		S3ObjectId: out.S3ObjectID,
	}), nil
}
