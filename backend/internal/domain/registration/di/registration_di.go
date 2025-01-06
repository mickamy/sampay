package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/registration/handler"
	"mickamy.com/sampay/internal/domain/registration/usecase"
)

type Repositories struct {
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet()

type UseCases struct {
	usecase.CreateAccount
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewCreateAccount,
)

type Handlers struct {
	*handler.Account
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewAccount,
)
