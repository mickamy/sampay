package di

import (
	"github.com/google/wire"

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
	CreateSession usecase.CreateSession
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewCreateSession,
)
