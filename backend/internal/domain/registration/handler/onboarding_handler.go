package handler

import (
	"context"
	"errors"

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
	"mickamy.com/sampay/internal/misc/i18n"
)

type Onboarding struct {
	getStep         usecase.GetOnboardingStep
	createPassword  usecase.CreatePassword
	createAttribute usecase.CreateUserAttribute
	createProfile   usecase.CreateUserProfile
	complete        usecase.CompleteOnboarding
}

func NewOnboarding(
	getStep usecase.GetOnboardingStep,
	createPassword usecase.CreatePassword,
	createAttribute usecase.CreateUserAttribute,
	createProfile usecase.CreateUserProfile,
	complete usecase.CompleteOnboarding,
) *Onboarding {
	return &Onboarding{
		getStep:         getStep,
		createPassword:  createPassword,
		createAttribute: createAttribute,
		createProfile:   createProfile,
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

func (h *Onboarding) CreateUserAttribute(
	ctx context.Context,
	req *connect.Request[registrationv1.CreateUserAttributeRequest],
) (*connect.Response[registrationv1.CreateUserAttributeResponse], error) {
	_, err := h.createAttribute.Do(ctx, usecase.CreateUserAttributeInput{
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
	res := connect.NewResponse(&registrationv1.CreateUserAttributeResponse{})
	return res, nil
}

func (h *Onboarding) CreateUserProfile(
	ctx context.Context,
	req *connect.Request[registrationv1.CreateUserProfileRequest],
) (*connect.Response[registrationv1.CreateUserProfileResponse], error) {
	lang := contexts.MustLanguage(ctx)
	obj, err := commonRequest.NewS3Object(req.Msg.Image)
	if err != nil {
		return nil, commonResponse.NewBadRequest(errors.New("invalid s3 object")).
			WithFieldViolation("image", i18n.MustLocalizeMessage(lang, i18n.Config{MessageID: i18n.CommonRequestErrorInvalid_s3_object})).
			AsConnectError()
	}

	_, err = h.createProfile.Do(ctx, usecase.CreateUserProfileInput{
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
	res := connect.NewResponse(&registrationv1.CreateUserProfileResponse{})
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
