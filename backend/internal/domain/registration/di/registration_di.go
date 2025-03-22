package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/registration/handler"
	"mickamy.com/sampay/internal/domain/registration/repository"
	"mickamy.com/sampay/internal/domain/registration/usecase"
)

type Repositories struct {
	repository.UsageCategory
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet(
	repository.NewUsageCategory,
)

type UseCases struct {
	usecase.CreatePassword
	usecase.CompleteOnboarding
	usecase.GetOnboardingStep
	usecase.ListUsageCategories
	usecase.UpdateUserAttribute
	usecase.UpdateUserLinks
	usecase.UpdateUserProfile
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewCreatePassword,
	usecase.NewCompleteOnboarding,
	usecase.NewGetOnboardingStep,
	usecase.NewListUsageCategories,
	usecase.NewUpdateUserAttribute,
	usecase.NewUpdateUserLinks,
	usecase.NewUpdateUserProfile,
)

type Handlers struct {
	*handler.Onboarding
	*handler.UsageCategory
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewOnboarding,
	handler.NewUsageCategory,
)
