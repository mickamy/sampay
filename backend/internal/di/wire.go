//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
	auth "mickamy.com/sampay/internal/domain/auth/di"
	registration "mickamy.com/sampay/internal/domain/registration/di"
	user "mickamy.com/sampay/internal/domain/user/di"
)

func InitConfigs() Configs {
	return NewConfigs()
}

func InitInfras() (Infras, error) {
	wire.Build(infras, configSet, wire.Struct(new(Infras), "*"))
	return Infras{}, nil
}

func InitAuthRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) auth.Repositories {
	wire.Build(
		auth.RepositorySet,
		wire.Struct(new(auth.Repositories), "*"),
	)
	return auth.Repositories{}
}

func InitAuthUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) auth.UseCases {
	wire.Build(
		auth.UseCaseSet,
		auth.RepositorySet,
		user.RepositorySet,
		wire.Struct(new(auth.UseCases), "*"),
	)
	return auth.UseCases{}
}

func InitAuthHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) auth.Handlers {
	wire.Build(
		auth.HandlerSet,
		auth.UseCaseSet,
		auth.RepositorySet,
		user.RepositorySet,
		wire.Struct(new(auth.Handlers), "*"),
	)
	return auth.Handlers{}
}

func InitRegistrationUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) registration.UseCases {
	wire.Build(
		registration.UseCaseSet,
		auth.RepositorySet,
		user.RepositorySet,
		wire.Struct(new(registration.UseCases), "*"),
	)
	return registration.UseCases{}
}

func InitRegistrationHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) registration.Handlers {
	wire.Build(
		registration.HandlerSet,
		registration.UseCaseSet,
		auth.RepositorySet,
		user.RepositorySet,
		wire.Struct(new(registration.Handlers), "*"),
	)
	return registration.Handlers{}
}

func InitUserRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) user.Repositories {
	wire.Build(
		user.RepositorySet,
		wire.Struct(new(user.Repositories), "*"),
	)
	return user.Repositories{}
}
