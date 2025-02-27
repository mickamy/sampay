package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/notification/handler"
	"mickamy.com/sampay/internal/domain/notification/repository"
	"mickamy.com/sampay/internal/domain/notification/usecase"
)

type Repositories struct {
	repository.Notification
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet(
	repository.NewNotification,
)

type UseCases struct {
	usecase.CountUnreadNotifications
	usecase.ListNotifications
	usecase.ReadNotification
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewCountUnreadNotifications,
	usecase.NewListNotifications,
	usecase.NewReadNotification,
)

type Handlers struct {
	*handler.Notification
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewNotification,
)
