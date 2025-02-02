package handler

import (
	"context"
	"errors"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/user/v1/userv1connect"
	userv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/user/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	"mickamy.com/sampay/internal/domain/common/dto/request"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/domain/user/usecase"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/misc/i18n"
)

type UserProfile struct {
	update      usecase.UpdateUserProfile
	updateImage usecase.UpdateUserProfileImage
}

func NewUserProfile(
	update usecase.UpdateUserProfile,
	updateImage usecase.UpdateUserProfileImage,
) *UserProfile {
	return &UserProfile{
		update:      update,
		updateImage: updateImage,
	}
}

func (h UserProfile) UpdateUserProfile(ctx context.Context, req *connect.Request[userv1.UpdateUserProfileRequest]) (*connect.Response[userv1.UpdateUserProfileResponse], error) {
	_, err := h.update.Do(ctx, usecase.UpdateUserProfileInput{
		Name: req.Msg.Name,
		Slug: req.Msg.Slug,
		Bio:  req.Msg.Bio,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			if errors.Is(err, userModel.ErrUserSlugAlreadyTaken) {
				return nil, localizable.AsFieldViolations("slug").AsConnectError()
			}
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.UpdateUserProfileResponse{})
	return res, nil
}

func (h UserProfile) UpdateUserProfileImage(ctx context.Context, req *connect.Request[userv1.UpdateUserProfileImageRequest]) (*connect.Response[userv1.UpdateUserProfileImageResponse], error) {
	lang := contexts.MustLanguage(ctx)
	obj, err := request.NewS3Object(req.Msg.Image)
	if err != nil {
		return nil, commonResponse.NewBadRequest(errors.New("invalid s3 object")).
			WithFieldViolation("s3_object", i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.CommonRequestErrorInvalid_s3_object})).
			AsConnectError()
	}

	_, err = h.updateImage.Do(ctx, usecase.UpdateUserProfileImageInput{
		Image: obj,
	})
	if err != nil {
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&userv1.UpdateUserProfileImageResponse{})
	return res, nil
}

var _ userv1connect.UserProfileServiceHandler = (*UserProfile)(nil)
