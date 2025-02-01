package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/message/repository"
	"mickamy.com/sampay/internal/domain/message/usecase"
)

type Repositories struct {
	repository.Message
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet(
	repository.NewMessage,
)

type UseCases struct {
	usecase.CreateMessage
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewCreateMessage,
)

type Handlers struct {
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet()
