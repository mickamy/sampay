package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/user/handler"
	"mickamy.com/sampay/internal/domain/user/repository"
	"mickamy.com/sampay/internal/domain/user/usecase"
)

type Repositories struct {
	repository.User
	repository.UserAttribute
	repository.UserLinkProvider
	repository.UserLink
	repository.UserProfile
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet(
	repository.NewUser,
	repository.NewUserAttribute,
	repository.NewUserLinkProvider,
	repository.NewUserLink,
	repository.NewUserProfile,
)

type UseCases struct {
	usecase.CreateUserLink
	usecase.DeleteUserLink
	usecase.GetMe
	usecase.GetUser
	usecase.ListUserLink
	usecase.UpdateUserLink
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewCreateUserLink,
	usecase.NewDeleteUserLink,
	usecase.NewGetMe,
	usecase.NewGetUser,
	usecase.NewListUserLink,
	usecase.NewUpdateUserLink,
)

type Handlers struct {
	*handler.UserLink
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewUserLink,
)
