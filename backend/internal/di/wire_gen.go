// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/cli/infra/storage/kvs"
	"mickamy.com/sampay/internal/domain/auth/di"
	"mickamy.com/sampay/internal/domain/auth/handler"
	"mickamy.com/sampay/internal/domain/auth/repository"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	di2 "mickamy.com/sampay/internal/domain/user/di"
	repository2 "mickamy.com/sampay/internal/domain/user/repository"
)

// Injectors from wire.go:

func InitInfras() (Infras, error) {
	databaseConfig := config.Database()
	db, err := provideDB(databaseConfig)
	if err != nil {
		return Infras{}, err
	}
	readWriter, err := provideReadWriter(databaseConfig)
	if err != nil {
		return Infras{}, err
	}
	writer, err := provideWriter(databaseConfig)
	if err != nil {
		return Infras{}, err
	}
	reader, err := provideReader(databaseConfig)
	if err != nil {
		return Infras{}, err
	}
	kvsConfig := config.KVS()
	v, err := provideKVS(kvsConfig)
	if err != nil {
		return Infras{}, err
	}
	diInfras := Infras{
		DB:         db,
		ReadWriter: readWriter,
		Writer:     writer,
		Reader:     reader,
		KVS:        v,
	}
	return diInfras, nil
}

func InitAuthRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di.Repositories {
	authentication := repository.NewAuthentication(db)
	session := repository.NewSession(kvs2)
	repositories := di.Repositories{
		Authentication: authentication,
		Session:        session,
	}
	return repositories
}

func InitAuthUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di.UseCases {
	authentication := repository.NewAuthentication(db)
	session := repository.NewSession(kvs2)
	user := repository2.NewUser(db)
	createSession := usecase.NewCreateSession(reader, authentication, session, user)
	refreshSession := usecase.NewRefreshSession(session)
	deleteSession := usecase.NewDeleteSession(session)
	useCases := di.UseCases{
		CreateSession:  createSession,
		RefreshSession: refreshSession,
		DeleteSession:  deleteSession,
	}
	return useCases
}

func InitAuthHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di.Handlers {
	authentication := repository.NewAuthentication(db)
	session := repository.NewSession(kvs2)
	user := repository2.NewUser(db)
	createSession := usecase.NewCreateSession(reader, authentication, session, user)
	refreshSession := usecase.NewRefreshSession(session)
	handlerSession := handler.NewSession(createSession, refreshSession)
	handlers := di.Handlers{
		Session: handlerSession,
	}
	return handlers
}

func InitUserRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di2.Repositories {
	user := repository2.NewUser(db)
	repositories := di2.Repositories{
		User: user,
	}
	return repositories
}

// wire.go:

func InitConfigs() Configs {
	return NewConfigs()
}
