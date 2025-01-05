package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/auth/handler"
	"mickamy.com/sampay/internal/domain/auth/repository"
	"mickamy.com/sampay/internal/domain/auth/usecase"
)

type Repositories struct {
	repository.Authentication
	repository.Session
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet(
	repository.NewAuthentication,
	repository.NewSession,
)

type UseCases struct {
	usecase.AuthenticateUser
	usecase.CreateSession
	usecase.RefreshSession
	usecase.DeleteSession
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewAuthenticateUser,
	usecase.NewCreateSession,
	usecase.NewRefreshSession,
	usecase.NewDeleteSession,
)

type Handlers struct {
	*handler.Session
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewSession,
)
