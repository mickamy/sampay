package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/auth/handler"
	"mickamy.com/sampay/internal/domain/auth/repository"
	"mickamy.com/sampay/internal/domain/auth/usecase"
)

type Repositories struct {
	repository.Authentication
	repository.EmailVerification
	repository.Session
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet(
	repository.NewAuthentication,
	repository.NewEmailVerification,
	repository.NewSession,
)

type UseCases struct {
	usecase.AuthenticateUser
	usecase.CreateSession
	usecase.DeleteSession
	usecase.RefreshSession
	usecase.RequestEmailVerification
	usecase.ResetPassword
	usecase.VerifyEmail
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewAuthenticateUser,
	usecase.NewCreateSession,
	usecase.NewDeleteSession,
	usecase.NewRefreshSession,
	usecase.NewResetPassword,
	usecase.NewRequestEmailVerification,
	usecase.NewVerifyEmail,
)

type Handlers struct {
	*handler.EmailVerification
	*handler.PasswordReset
	*handler.Session
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewEmailVerification,
	handler.NewPasswordReset,
	handler.NewSession,
)
