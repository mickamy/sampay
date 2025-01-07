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
	getStep usecase.GetOnboardingStep
}

func NewOnboarding(
	getStep usecase.GetOnboardingStep,
) *Onboarding {
	return &Onboarding{
		getStep: getStep,
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

func (h *Onboarding) PostUsageCategory(
	ctx context.Context,
	req *connect.Request[registrationv1.PostUsageCategoryRequest],
) (*connect.Response[registrationv1.PostUsageCategoryResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (h *Onboarding) PostUserProfile(
	ctx context.Context,
	req *connect.Request[registrationv1.PostUserProfileRequest],
) (*connect.Response[registrationv1.PostUserProfileResponse], error) {
	//TODO implement me
	panic("implement me")
}

var _ registrationv1connect.OnboardingServiceHandler = (*Onboarding)(nil)
