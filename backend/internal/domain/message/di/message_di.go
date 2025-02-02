package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/message/handler"
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
	usecase.SendMessage
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewSendMessage,
)

type Handlers struct {
	*handler.Message
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewMessage,
)
