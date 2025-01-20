package handler

import (
	"context"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/registration/v1/registrationv1connect"
	registrationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/registration/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/domain/registration/model"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/slices"
)

type UsageCategory struct {
	list usecase.ListUsageCategories
}

func NewUsageCategory(
	list usecase.ListUsageCategories,
) *UsageCategory {
	return &UsageCategory{
		list: list,
	}
}

func (h *UsageCategory) ListUsageCategories(
	ctx context.Context,
	req *connect.Request[registrationv1.ListUsageCategoriesRequest],
) (*connect.Response[registrationv1.ListUsageCategoriesResponse], error) {
	out, err := h.list.Do(ctx, usecase.ListUsageCategoriesInput{})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.ListUsageCategoriesResponse{
		Categories: slices.Map(out.Categories, func(c model.UsageCategory) *registrationv1.UsageCategory {
			return &registrationv1.UsageCategory{
				Type:         c.Type,
				DisplayOrder: int32(c.DisplayOrder),
			}
		}),
	})
	return res, nil
}

var _ registrationv1connect.UsageCategoryServiceHandler = (*UsageCategory)(nil)
