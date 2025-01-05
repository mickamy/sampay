//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/cli/infra/kvs"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	auth "mickamy.com/sampay/internal/domain/auth/di"
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

func InitUserRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) user.Repositories {
	wire.Build(
		user.RepositorySet,
		wire.Struct(new(user.Repositories), "*"),
	)
	return user.Repositories{}
}
