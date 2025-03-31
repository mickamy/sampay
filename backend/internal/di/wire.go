//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
	auth "mickamy.com/sampay/internal/domain/auth/di"
	common "mickamy.com/sampay/internal/domain/common/di"
	message "mickamy.com/sampay/internal/domain/message/di"
	notification "mickamy.com/sampay/internal/domain/notification/di"
	oauth "mickamy.com/sampay/internal/domain/oauth/di"
	registration "mickamy.com/sampay/internal/domain/registration/di"
	user "mickamy.com/sampay/internal/domain/user/di"
	"mickamy.com/sampay/internal/infra/storage/database"
	"mickamy.com/sampay/internal/job"
)

func InitInfras() (Infras, error) {
	wire.Build(infraSet, configSet, wire.Struct(new(Infras), "*"))
	return Infras{}, nil
}

func InitLibs() Libs {
	wire.Build(libSet, configSet, wire.Struct(new(Libs), "*"))
	return Libs{}
}

func InitJobs() job.Jobs {
	wire.Build(
		jobSet,
		configSet,
		libSet,
		wire.Struct(new(job.Jobs), "*"),
	)
	return job.Jobs{}
}

func InitProducers() Producers {
	wire.Build(
		producerSet,
		configSet,
		wire.Struct(new(Producers), "*"),
	)
	return Producers{}
}

func InitConsumers() Consumers {
	wire.Build(
		consumerSet,
		configSet,
		InitJobs,
		wire.Struct(new(Consumers), "*"),
	)
	return Consumers{}
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
		configSet,
		producerSet,
		auth.RepositorySet,
		user.RepositorySet,
		wire.Struct(new(auth.UseCases), "*"),
	)
	return auth.UseCases{}
}

func InitAuthHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) auth.Handlers {
	wire.Build(
		auth.HandlerSet,
		configSet,
		producerSet,
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

func InitMessageRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) message.Repositories {
	wire.Build(
		message.RepositorySet,
		wire.Struct(new(message.Repositories), "*"),
	)
	return message.Repositories{}
}

func InitMessageUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) message.UseCases {
	wire.Build(
		message.UseCaseSet,
		configSet,
		producerSet,
		notification.RepositorySet,
		user.RepositorySet,
		message.RepositorySet,
		wire.Struct(new(message.UseCases), "*"),
	)
	return message.UseCases{}
}

func InitMessageHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) message.Handlers {
	wire.Build(
		message.HandlerSet,
		configSet,
		producerSet,
		notification.RepositorySet,
		user.RepositorySet,
		message.RepositorySet,
		message.UseCaseSet,
		wire.Struct(new(message.Handlers), "*"),
	)
	return message.Handlers{}
}

func InitNotificationRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) notification.Repositories {
	wire.Build(
		notification.RepositorySet,
		wire.Struct(new(notification.Repositories), "*"),
	)
	return notification.Repositories{}
}

func InitNotificationUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) notification.UseCases {
	wire.Build(
		notification.UseCaseSet,
		notification.RepositorySet,
		wire.Struct(new(notification.UseCases), "*"),
	)
	return notification.UseCases{}
}

func InitNotificationHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) notification.Handlers {
	wire.Build(
		notification.HandlerSet,
		notification.RepositorySet,
		notification.UseCaseSet,
		wire.Struct(new(notification.Handlers), "*"),
	)
	return notification.Handlers{}
}

func InitOAuthUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) oauth.UseCases {
	wire.Build(
		oauth.UseCaseSet,
		configSet,
		libSet,
		auth.RepositorySet,
		user.RepositorySet,
		wire.Struct(new(oauth.UseCases), "*"),
	)
	return oauth.UseCases{}
}

func InitOAuthHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) oauth.Handlers {
	wire.Build(
		oauth.HandlerSet,
		configSet,
		libSet,
		auth.RepositorySet,
		user.RepositorySet,
		oauth.UseCaseSet,
		wire.Struct(new(oauth.Handlers), "*"),
	)
	return oauth.Handlers{}
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
		configSet,
		libSet,
		common.RepositorySet,
		user.RepositorySet,
		wire.Struct(new(user.UseCases), "*"),
	)
	return user.UseCases{}
}

func InitUserHandler(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *kvs.KVS) user.Handlers {
	wire.Build(
		user.HandlerSet,
		configSet,
		libSet,
		common.RepositorySet,
		user.RepositorySet,
		user.UseCaseSet,
		wire.Struct(new(user.Handlers), "*"),
	)
	return user.Handlers{}
}
