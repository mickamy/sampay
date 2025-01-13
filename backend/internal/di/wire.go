//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
	auth "mickamy.com/sampay/internal/domain/auth/di"
	common "mickamy.com/sampay/internal/domain/common/di"
	registration "mickamy.com/sampay/internal/domain/registration/di"
	user "mickamy.com/sampay/internal/domain/user/di"
)

func InitInfras() (Infras, error) {
	wire.Build(infraSet, configSet, wire.Struct(new(Infras), "*"))
	return Infras{}, nil
}

func InitLibs() Libs {
	wire.Build(libSet, configSet, wire.Struct(new(Libs), "*"))
	return Libs{}
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

func InitCommonRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) common.Repositories {
	wire.Build(
		common.RepositorySet,
		wire.Struct(new(common.Repositories), "*"),
	)
	return common.Repositories{}
}

func InitCommonUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) common.UseCases {
	wire.Build(
		common.UseCaseSet,
		configSet,
		libSet,
		wire.Struct(new(common.UseCases), "*"),
	)
	return common.UseCases{}
}

func InitCommonHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) common.Handlers {
	wire.Build(
		common.HandlerSet,
		configSet,
		libSet,
		common.UseCaseSet,
		wire.Struct(new(common.Handlers), "*"),
	)
	return common.Handlers{}
}

func InitRegistrationRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) registration.Repositories {
	wire.Build(
		registration.RepositorySet,
		wire.Struct(new(registration.Repositories), "*"),
	)
	return registration.Repositories{}
}

func InitRegistrationUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) registration.UseCases {
	wire.Build(
		registration.UseCaseSet,
		auth.RepositorySet,
		user.RepositorySet,
		registration.RepositorySet,
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
		registration.RepositorySet,
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

func InitUserUseCase(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) user.UseCases {
	wire.Build(
		user.UseCaseSet,
		common.RepositorySet,
		user.RepositorySet,
		wire.Struct(new(user.UseCases), "*"),
	)
	return user.UseCases{}
}

func InitUserHandler(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) user.Handlers {
	wire.Build(
		user.HandlerSet,
		user.RepositorySet,
		user.UseCaseSet,
		wire.Struct(new(user.Handlers), "*"),
	)
	return user.Handlers{}
}
