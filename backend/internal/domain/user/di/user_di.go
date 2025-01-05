package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/user/repository"
)

type Repositories struct {
	repository.User
}

type UseCases struct {
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet(
	repository.NewUser,
)
