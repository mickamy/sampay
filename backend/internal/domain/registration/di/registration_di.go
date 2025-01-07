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
	usecase.CreateAccount
	usecase.GetOnboardingStep
	usecase.ListUsageCategories
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewCreateAccount,
	usecase.NewGetOnboardingStep,
	usecase.NewListUsageCategories,
)

type Handlers struct {
	*handler.Account
	*handler.Onboarding
	*handler.UsageCategory
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewAccount,
	handler.NewOnboarding,
	handler.NewUsageCategory,
)
