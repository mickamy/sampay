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
	"mickamy.com/sampay/internal/lib/contexts"
)

type Onboarding struct {
	getStep         usecase.GetOnboardingStep
	createPassword  usecase.CreatePassword
	createAttribute usecase.CreateUserAttribute
	createProfile   usecase.CreateUserProfile
}

func NewOnboarding(
	getStep usecase.GetOnboardingStep,
	createPassword usecase.CreatePassword,
	createAttribute usecase.CreateUserAttribute,
	createProfile usecase.CreateUserProfile,
) *Onboarding {
	return &Onboarding{
		getStep:         getStep,
		createPassword:  createPassword,
		createAttribute: createAttribute,
		createProfile:   createProfile,
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
	_, err := h.createProfile.Do(ctx, usecase.CreateUserProfileInput{
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
	res := connect.NewResponse(&registrationv1.CreateUserProfileResponse{})
	return res, nil
}

var _ registrationv1connect.OnboardingServiceHandler = (*Onboarding)(nil)
