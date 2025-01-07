package handler

import (
	"context"

	"buf.build/gen/go/mickamy/sampay/connectrpc/go/registration/v1/registrationv1connect"
	registrationv1 "buf.build/gen/go/mickamy/sampay/protocolbuffers/go/registration/v1"
	"connectrpc.com/connect"
	"github.com/mickamy/slogger"

	dto "mickamy.com/sampay/internal/domain/common/dto"
	"mickamy.com/sampay/internal/domain/registration/usecase"
)

type Onboarding struct {
	getStep         usecase.GetOnboardingStep
	createAttribute usecase.CreateUserAttribute
	createProfile   usecase.CreateUserProfile
}

func NewOnboarding(
	getStep usecase.GetOnboardingStep,
	createAttribute usecase.CreateUserAttribute,
	createProfile usecase.CreateUserProfile,
) *Onboarding {
	return &Onboarding{
		getStep:         getStep,
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
		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.GetOnboardingStepResponse{
		Step: out.Step.String(),
	})
	return res, nil
}

func (h *Onboarding) CreateUserAttribute(
	ctx context.Context,
	req *connect.Request[registrationv1.CreateUserAttributeRequest],
) (*connect.Response[registrationv1.CreateUserAttributeResponse], error) {
	_, err := h.createAttribute.Do(ctx, usecase.CreateUserAttributeInput{})
	if err != nil {
		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.CreateUserAttributeResponse{})
	return res, nil
}

func (h *Onboarding) CreateUserProfile(
	ctx context.Context,
	req *connect.Request[registrationv1.CreateUserProfileRequest],
) (*connect.Response[registrationv1.CreateUserProfileResponse], error) {
	_, err := h.createProfile.Do(ctx, usecase.CreateUserProfileInput{
		Name: req.Msg.Name,
		Bio:  req.Msg.Bio,
	})
	if err != nil {
		slogger.ErrorCtx(ctx, "failed to execute use case", "err", err)
		return nil, dto.NewInternalError(ctx, err).AsConnectError()
	}
	res := connect.NewResponse(&registrationv1.CreateUserProfileResponse{})
	return res, nil
}

var _ registrationv1connect.OnboardingServiceHandler = (*Onboarding)(nil)
