package handler

import (
	"context"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/user/v1/userv1connect"
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	dto "mickamy.com/sampay/internal/domain/common/dto"
	"mickamy.com/sampay/internal/domain/user/dto/response"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/ptr"
)

type UserLink struct {
	create usecase.CreateUserLink
	list   usecase.ListUserLink
	update usecase.UpdateUserLink
	delete usecase.DeleteUserLink
}

func NewUserLink(
	create usecase.CreateUserLink,
	list usecase.ListUserLink,
	update usecase.UpdateUserLink,
	delete usecase.DeleteUserLink,
) *UserLink {
	return &UserLink{
		create: create,
		list:   list,
		update: update,
		delete: delete,
	}
}

func (h UserLink) ListUserLink(ctx context.Context, req *connect.Request[userv1.ListUserLinkRequest]) (*connect.Response[userv1.ListUserLinkResponse], error) {
	out, err := h.list.Do(ctx, usecase.ListUserLinkInput{
		UserID: req.Msg.UserId,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := dto.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.ListUserLinkResponse{
		Links: response.NewUserLinks(out.Links),
	})
	return res, nil
}

func (h UserLink) CreateUserLink(ctx context.Context, req *connect.Request[userv1.CreateUserLinkRequest]) (*connect.Response[userv1.CreateUserLinkResponse], error) {
	_, err := h.create.Do(ctx, usecase.CreateUserLinkInput{
		ProviderType: userModel.MustNewLinkProviderType(req.Msg.ProviderType),
		URI:          req.Msg.Uri,
		Name:         req.Msg.Name,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := dto.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.CreateUserLinkResponse{})
	return res, nil
}

func (h UserLink) UpdateUserLink(ctx context.Context, req *connect.Request[userv1.UpdateUserLinkRequest]) (*connect.Response[userv1.UpdateUserLinkResponse], error) {
	_, err := h.update.Do(ctx, usecase.UpdateUserLinkInput{
		ProviderType: ptr.Map(req.Msg.ProviderType, func(v *string) *userModel.UserLinkProviderType {
			return ptr.Of(userModel.MustNewLinkProviderType(*v))
		}),
		URI:  req.Msg.Uri,
		Name: req.Msg.Name,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := dto.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.UpdateUserLinkResponse{})
	return res, nil
}

func (h UserLink) DeleteUserLink(ctx context.Context, c *connect.Request[userv1.DeleteUserLinkRequest]) (*connect.Response[userv1.DeleteUserLinkResponse], error) {
	_, err := h.delete.Do(ctx, usecase.DeleteUserLinkInput{
		ID: c.Msg.Id,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := dto.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.DeleteUserLinkResponse{})
	return res, nil
}

var _ userv1connect.UserLinkServiceHandler = (*UserLink)(nil)
