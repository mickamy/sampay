package handler

import (
	"context"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/user/v1/userv1connect"
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	commonRequest "mickamy.com/sampay/internal/domain/common/dto/request"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/domain/user/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
)

type UserProfile struct {
	update      usecase.UpdateUserProfile
	deleteImage usecase.DeleteUserProfileImage
}

func NewUserProfile(
	update usecase.UpdateUserProfile,
	deleteImage usecase.DeleteUserProfileImage,
) *UserProfile {
	return &UserProfile{
		update:      update,
		deleteImage: deleteImage,
	}
}

func (h UserProfile) UpdateUserProfile(ctx context.Context, req *connect.Request[userv1.UpdateUserProfileRequest]) (*connect.Response[userv1.UpdateUserProfileResponse], error) {
	_, err := h.update.Do(ctx, usecase.UpdateUserProfileInput{
		Name:  req.Msg.Name,
		Bio:   req.Msg.Bio,
		Image: commonRequest.NewS3Object(req.Msg.Image),
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.UpdateUserProfileResponse{})
	return res, nil
}

func (h UserProfile) DeleteUserProfileImage(ctx context.Context, req *connect.Request[userv1.DeleteUserProfileImageRequest]) (*connect.Response[userv1.DeleteUserProfileImageResponse], error) {
	_, err := h.update.Do(ctx, usecase.UpdateUserProfileInput{})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.DeleteUserProfileImageResponse{})
	return res, nil
}

var _ userv1connect.UserProfileServiceHandler = (*UserProfile)(nil)
