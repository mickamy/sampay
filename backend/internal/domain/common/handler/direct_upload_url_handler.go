package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/common/v1/commonv1connect"
	"buf.build/gen/go/mickamy/sampay/protocolbuffers/go/common/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	dto "mickamy.com/sampay/internal/domain/common/dto"
	"mickamy.com/sampay/internal/domain/common/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/misc/i18n"
)

type DirectUploadURL struct {
	create usecase.CreateDirectUploadURL
}

func NewDirectUploadURL(
	create usecase.CreateDirectUploadURL,
) *DirectUploadURL {
	return &DirectUploadURL{
		create: create,
	}
}

func (h *DirectUploadURL) Request(
	ctx context.Context,
	req *connect.Request[commonv1.DirectUploadURLRequest],
) (*connect.Response[commonv1.DirectUploadURLResponse], error) {
	lang := contexts.MustLanguage(ctx)
	obj := dto.NewS3Object(req.Msg.S3Object)
	if obj == nil {
		return nil, dto.NewBadRequest(errors.New("invalid s3 object")).
			WithMessage(i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: "common.handler.direct_upload_url.error.invalid_s3_object"})).
			AsConnectError()
	}
	out, err := h.create.Do(ctx, usecase.CreateDirectUploadURLInput{
		S3Object: *obj,
	})
	if err != nil {
		if localizable := dto.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&commonv1.DirectUploadURLResponse{
		Url: out.URL,
	})
	return res, nil
}

var _ commonv1connect.DirectUploadURLServiceHandler = (*DirectUploadURL)(nil)
