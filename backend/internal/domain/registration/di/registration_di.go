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
	usecase.CreatePassword
	usecase.CreateUserAttribute
	usecase.CreateUserProfile
	usecase.GetOnboardingStep
	usecase.ListUsageCategories
	usecase.RequestEmailVerification
	usecase.VerifyEmail
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewCreateAccount,
	usecase.NewCreatePassword,
	usecase.NewCreateUserAttribute,
	usecase.NewCreateUserProfile,
	usecase.NewGetOnboardingStep,
	usecase.NewListUsageCategories,
	usecase.NewRequestEmailVerification,
	usecase.NewVerifyEmail,
)

type Handlers struct {
	*handler.Account
	*handler.EmailVerification
	*handler.Onboarding
	*handler.UsageCategory
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewAccount,
	handler.NewEmailVerification,
	handler.NewOnboarding,
	handler.NewUsageCategory,
)
