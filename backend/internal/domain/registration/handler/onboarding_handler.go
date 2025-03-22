package handler

import (
	"context"
	"errors"
	"strconv"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/registration/v1/registrationv1connect"
	registrationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/registration/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	authResponse "mickamy.com/sampay/internal/domain/auth/dto/response"
	commonRequest "mickamy.com/sampay/internal/domain/common/dto/request"
	commonResponse "mickamy.com/sampay/internal/domain/common/dto/response"
	"mickamy.com/sampay/internal/domain/registration/usecase"
	userModel "mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/ptr"
	"mickamy.com/sampay/internal/misc/i18n"
)

type Onboarding struct {
	getStep         usecase.GetOnboardingStep
	createPassword  usecase.CreatePassword
	updateAttribute usecase.UpdateUserAttribute
	updateProfile   usecase.UpdateUserProfile
	updateLinks     usecase.UpdateUserLinks
	complete        usecase.CompleteOnboarding
}

func NewOnboarding(
	getStep usecase.GetOnboardingStep,
	createPassword usecase.CreatePassword,
	updateAttribute usecase.UpdateUserAttribute,
	updateProfile usecase.UpdateUserProfile,
	updateLinks usecase.UpdateUserLinks,
	complete usecase.CompleteOnboarding,
) *Onboarding {
	return &Onboarding{
		getStep:         getStep,
		createPassword:  createPassword,
		updateAttribute: updateAttribute,
		updateProfile:   updateProfile,
		updateLinks:     updateLinks,
		complete:        complete,
	}
}

func (h *Onboarding) GetOnboardingStep(
	ctx context.Context,
	req *connect.Request[registrationv1.GetOnboardingStepRequest],
) (*connect.Response[registrationv1.GetOnboardingStepResponse], error) {
	out, err := h.getStep.Do(ctx, usecase.GetOnboardingStepInput{})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.GetOnboardingStepResponse{
		Step: out.Step.String(),
	})
	return res, nil
}

func (h *Onboarding) CreatePassword(
	ctx context.Context,
	req *connect.Request[registrationv1.CreatePasswordRequest],
) (*connect.Response[registrationv1.CreatePasswordResponse], error) {
	got, err := h.createPassword.Do(ctx, usecase.CreatePasswordInput{
		Password: req.Msg.Password,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			if errors.Is(err, usecase.ErrCreatePasswordEmailVerificationInvalidToken) || errors.Is(err, usecase.ErrCreatePasswordEmailVerificationAlreadyConsumed) {
				return nil, localizable.AsFieldViolations("token").AsConnectError()
			}

			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.CreatePasswordResponse{
		Tokens: authResponse.NewTokens(got.Session.Tokens),
	})
	return res, nil
}

func (h *Onboarding) UpdateUserAttribute(
	ctx context.Context,
	req *connect.Request[registrationv1.UpdateUserAttributeRequest],
) (*connect.Response[registrationv1.UpdateUserAttributeResponse], error) {
	_, err := h.updateAttribute.Do(ctx, usecase.UpdateUserAttributeInput{
		UsageCategoryType: req.Msg.CategoryType,
	})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.UpdateUserAttributeResponse{})
	return res, nil
}

func (h *Onboarding) UpdateUserProfile(
	ctx context.Context,
	req *connect.Request[registrationv1.UpdateUserProfileRequest],
) (*connect.Response[registrationv1.UpdateUserProfileResponse], error) {
	lang := contexts.MustLanguage(ctx)
	obj, err := commonRequest.NewS3Object(req.Msg.Image)
	if err != nil {
		return nil, commonResponse.NewBadRequest(errors.New("invalid s3 object")).
			WithFieldViolation("image", i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.CommonRequestErrorInvalid_s3_object})).
			AsConnectError()
	}

	_, err = h.updateProfile.Do(ctx, usecase.UpdateUserProfileInput{
		Name:  req.Msg.Name,
		Slug:  req.Msg.Slug,
		Bio:   req.Msg.Bio,
		Image: obj,
	})
	if err != nil {
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			if errors.Is(err, userModel.ErrUserSlugAlreadyTaken) {
				return nil, localizable.AsFieldViolations("slug").AsConnectError()
			}
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.UpdateUserProfileResponse{})
	return res, nil
}

type userLinkS3ObjectError struct {
	index int
	err   error
}

func (h *Onboarding) UpdateUserLinks(
	ctx context.Context,
	req *connect.Request[registrationv1.UpdateUserLinksRequest],
) (*connect.Response[registrationv1.UpdateUserLinksResponse], error) {
	var errs []userLinkS3ObjectError
	links := make([]usecase.UserLink, 0, len(req.Msg.Links))
	for i, link := range req.Msg.Links {
		obj, err := commonRequest.NewS3Object(link.QrCode)
		if err != nil {
			errs = append(errs, userLinkS3ObjectError{
				index: i,
				err:   err,
			})
		}
		links = append(links, usecase.UserLink{
			ID:           ptr.ZeroIfNull(link.Id),
			ProviderType: userModel.UserLinkProviderType(link.ProviderType),
			URI:          link.Uri,
			Name:         link.Name,
			QRCode:       obj,
		})
	}

	if len(errs) > 0 {
		lang := contexts.MustLanguage(ctx)
		res := commonResponse.NewBadRequest(errors.New("invalid s3 object"))
		for _, e := range errs {
			res = res.WithFieldViolation(
				"links["+strconv.Itoa(e.index)+"].qr_code",
				i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.CommonRequestErrorInvalid_s3_object}),
			)
		}
		return nil, res.AsConnectError()
	}

	_, err := h.updateLinks.Do(ctx, usecase.UpdateUserLinksInput{
		UserLinks: links,
	})

	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.UpdateUserLinksResponse{})
	return res, nil
}

func (h *Onboarding) CompleteOnboarding(
	ctx context.Context,
	req *connect.Request[registrationv1.CompleteOnboardingRequest],
) (*connect.Response[registrationv1.CompleteOnboardingResponse], error) {
	_, err := h.complete.Do(ctx, usecase.CompleteOnboardingInput{})
	if err != nil {
		lang := contexts.MustLanguage(ctx)
		if localizable := commonResponse.ParseLocalizableError(lang, err); localizable != nil {
			return nil, localizable.AsConnectError()
		}

		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, commonResponse.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.CompleteOnboardingResponse{})
	return res, nil
}

var _ registrationv1connect.OnboardingServiceHandler = (*Onboarding)(nil)
