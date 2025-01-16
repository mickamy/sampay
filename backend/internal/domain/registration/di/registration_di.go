package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/registration/handler"
	"mickamy.com/sampay/internal/domain/registration/repository"
	"mickamy.com/sampay/internal/domain/registration/usecase"
)

type Repositories struct {
	repository.EmailVerification
	repository.UsageCategory
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet(
	repository.NewEmailVerification,
	repository.NewUsageCategory,
)

type UseCases struct {
	usecase.CreateAccount
	usecase.CreateUserAttribute
	usecase.CreateUserProfile
	usecase.GetOnboardingStep
	usecase.ListUsageCategories
	usecase.RequestEmailVerification
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewCreateAccount,
	usecase.NewCreateUserAttribute,
	usecase.NewCreateUserProfile,
	usecase.NewGetOnboardingStep,
	usecase.NewListUsageCategories,
	usecase.NewRequestEmailVerification,
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
