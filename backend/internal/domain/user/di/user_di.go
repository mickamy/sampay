package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/user/repository"
)

type Repositories struct {
	repository.User
	repository.UserAttribute
	repository.UserProfile
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet(
	repository.NewUser,
	repository.NewUserAttribute,
	repository.NewUserProfile,
)

type UseCases struct {
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet()
