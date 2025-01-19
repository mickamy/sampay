package handler

import (
	"context"

	"buf.build/gen/go/mickamy/sampay/bufbuild/connect-go/user/v1/userv1connect"
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	commonRequest "mickamy.com/sampay/internal/domain/common/dto/request"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/domain/user/dto/response"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/misc/i18n"
)

type UserLink struct {
	create       usecase.CreateUserLink
	list         usecase.ListUserLink
	update       usecase.UpdateUserLink
	updateQRCode usecase.UpdateUserLinkQRCode
	delete       usecase.DeleteUserLink
}

func NewUserLink(
	create usecase.CreateUserLink,
	list usecase.ListUserLink,
	update usecase.UpdateUserLink,
	updateQRCode usecase.UpdateUserLinkQRCode,
	delete usecase.DeleteUserLink,
) *UserLink {
	return &UserLink{
		create:       create,
		list:         list,
		update:       update,
		updateQRCode: updateQRCode,
		delete:       delete,
	}
}

func (h UserLink) CreateUserLink(ctx context.Context, req *connect.Request[userv1.CreateUserLinkRequest]) (*connect.Response[userv1.CreateUserLinkResponse], error) {
	lang := contexts.MustLanguage(ctx)
	providerType, err := userModel.NewLinkProviderType(req.Msg.ProviderType)
	if err != nil {
		return nil, commonResponse.NewBadRequest(err, commonResponse.FieldViolation{
			Field:        "provider_type",
			Descriptions: []string{i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.UserHandlerUser_linkErrorInvalid_provider_type})},
		}).AsConnectError()
	}

	_, err = h.create.Do(ctx, usecase.CreateUserLinkInput{
		ProviderType: providerType,
		URI:          req.Msg.Uri,
		Name:         req.Msg.Name,
	})
	if err != nil {
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.CreateUserLinkResponse{})
	return res, nil
}

func (h UserLink) ListUserLink(ctx context.Context, req *connect.Request[userv1.ListUserLinkRequest]) (*connect.Response[userv1.ListUserLinkResponse], error) {
	out, err := h.list.Do(ctx, usecase.ListUserLinkInput{
		UserID: req.Msg.UserId,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.ListUserLinkResponse{
		Links: response.NewUserLinks(out.Links),
	})
	return res, nil
}

func (h UserLink) UpdateUserLink(ctx context.Context, req *connect.Request[userv1.UpdateUserLinkRequest]) (*connect.Response[userv1.UpdateUserLinkResponse], error) {
	lang := contexts.MustLanguage(ctx)
	var providerType *userModel.UserLinkProviderType
	if req.Msg.ProviderType != nil {
		pt, err := userModel.NewLinkProviderType(*req.Msg.ProviderType)
		if err != nil {
			return nil, commonResponse.NewBadRequest(err, commonResponse.FieldViolation{
				Field:        "provider_type",
				Descriptions: []string{i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.UserHandlerUser_linkErrorInvalid_provider_type})},
			}).AsConnectError()
		}
		providerType = &pt
	}

	_, err := h.update.Do(ctx, usecase.UpdateUserLinkInput{
		ID:           req.Msg.Id,
		ProviderType: providerType,
		URI:          req.Msg.Uri,
		Name:         req.Msg.Name,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.UpdateUserLinkResponse{})
	return res, nil
}

func (h UserLink) UpdateUserLinkQRCode(ctx context.Context, req *connect.Request[userv1.UpdateUserLinkQRCodeRequest]) (*connect.Response[userv1.UpdateUserLinkQRCodeResponse], error) {
	_, err := h.updateQRCode.Do(ctx, usecase.UpdateUserLinkQRCodeInput{
		ID:     req.Msg.Id,
		QRCode: commonRequest.NewS3Object(req.Msg.QrCode),
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.UpdateUserLinkQRCodeResponse{})
	return res, nil
}

func (h UserLink) DeleteUserLink(ctx context.Context, req *connect.Request[userv1.DeleteUserLinkRequest]) (*connect.Response[userv1.DeleteUserLinkResponse], error) {
	_, err := h.delete.Do(ctx, usecase.DeleteUserLinkInput{
		ID: req.Msg.Id,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.DeleteUserLinkResponse{})
	return res, nil
}

var _ userv1connect.UserLinkServiceHandler = (*UserLink)(nil)
