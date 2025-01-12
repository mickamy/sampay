// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/redis/go-redis/v9"
	"mickamy.com/sampay/config"
	"mickamy.com/sampay/internal/cli/infra/storage/database"
	"mickamy.com/sampay/internal/domain/auth/di"
	"mickamy.com/sampay/internal/domain/auth/handler"
	"mickamy.com/sampay/internal/domain/auth/repository"
	"mickamy.com/sampay/internal/domain/auth/usecase"
	di2 "mickamy.com/sampay/internal/domain/common/di"
	handler2 "mickamy.com/sampay/internal/domain/common/handler"
	usecase2 "mickamy.com/sampay/internal/domain/common/usecase"
	di3 "mickamy.com/sampay/internal/domain/registration/di"
	handler3 "mickamy.com/sampay/internal/domain/registration/handler"
	repository3 "mickamy.com/sampay/internal/domain/registration/repository"
	usecase3 "mickamy.com/sampay/internal/domain/registration/usecase"
	di4 "mickamy.com/sampay/internal/domain/user/di"
	handler4 "mickamy.com/sampay/internal/domain/user/handler"
	repository2 "mickamy.com/sampay/internal/domain/user/repository"
	usecase4 "mickamy.com/sampay/internal/domain/user/usecase"
	"mickamy.com/sampay/internal/lib/aws/s3"
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
	client, err := provideKVS(kvsConfig)
	if err != nil {
		return Infras{}, err
	}
	infras := Infras{
		DB:         db,
		ReadWriter: readWriter,
		Writer:     writer,
		Reader:     reader,
		KVS:        client,
	}
	return infras, nil
}

func InitLibs() Libs {
	awsConfig := config.AWS()
	client := s3.New(awsConfig)
	libs := Libs{
		Client: client,
	}
	return libs
}

func InitAuthRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di.Repositories {
	authentication := repository.NewAuthentication(db)
	session := repository.NewSession(kvs)
	repositories := di.Repositories{
		Authentication: authentication,
		Session:        session,
	}
	return repositories
}

func InitAuthUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di.UseCases {
	session := repository.NewSession(kvs)
	user := repository2.NewUser(db)
	authenticateUser := usecase.NewAuthenticateUser(reader, session, user)
	authentication := repository.NewAuthentication(db)
	createSession := usecase.NewCreateSession(reader, authentication, session, user)
	refreshSession := usecase.NewRefreshSession(session)
	deleteSession := usecase.NewDeleteSession(session)
	useCases := di.UseCases{
		AuthenticateUser: authenticateUser,
		CreateSession:    createSession,
		RefreshSession:   refreshSession,
		DeleteSession:    deleteSession,
	}
	return useCases
}

func InitAuthHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di.Handlers {
	authentication := repository.NewAuthentication(db)
	session := repository.NewSession(kvs)
	user := repository2.NewUser(db)
	createSession := usecase.NewCreateSession(reader, authentication, session, user)
	refreshSession := usecase.NewRefreshSession(session)
	deleteSession := usecase.NewDeleteSession(session)
	handlerSession := handler.NewSession(createSession, refreshSession, deleteSession)
	handlers := di.Handlers{
		Session: handlerSession,
	}
	return handlers
}

func InitCommonUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di2.UseCases {
	awsConfig := config.AWS()
	client := s3.New(awsConfig)
	createDirectUploadURL := usecase2.NewCreateDirectUploadURL(client)
	useCases := di2.UseCases{
		CreateDirectUploadURL: createDirectUploadURL,
	}
	return useCases
}

func InitCommonHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di2.Handlers {
	awsConfig := config.AWS()
	client := s3.New(awsConfig)
	createDirectUploadURL := usecase2.NewCreateDirectUploadURL(client)
	directUploadURL := handler2.NewDirectUploadURL(createDirectUploadURL)
	handlers := di2.Handlers{
		DirectUploadURL: directUploadURL,
	}
	return handlers
}

func InitRegistrationRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di3.Repositories {
	usageCategory := repository3.NewUsageCategory(db)
	repositories := di3.Repositories{
		UsageCategory: usageCategory,
	}
	return repositories
}

func InitRegistrationUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di3.UseCases {
	authentication := repository.NewAuthentication(db)
	session := repository.NewSession(kvs)
	user := repository2.NewUser(db)
	createAccount := usecase3.NewCreateAccount(writer, authentication, session, user)
	userAttribute := repository2.NewUserAttribute(db)
	createUserAttribute := usecase3.NewCreateUserAttribute(writer, userAttribute)
	userProfile := repository2.NewUserProfile(db)
	createUserProfile := usecase3.NewCreateUserProfile(writer, userProfile)
	getOnboardingStep := usecase3.NewGetOnboardingStep(reader, user)
	usageCategory := repository3.NewUsageCategory(db)
	listUsageCategories := usecase3.NewListUsageCategories(reader, usageCategory)
	useCases := di3.UseCases{
		CreateAccount:       createAccount,
		CreateUserAttribute: createUserAttribute,
		CreateUserProfile:   createUserProfile,
		GetOnboardingStep:   getOnboardingStep,
		ListUsageCategories: listUsageCategories,
	}
	return useCases
}

func InitRegistrationHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di3.Handlers {
	authentication := repository.NewAuthentication(db)
	session := repository.NewSession(kvs)
	user := repository2.NewUser(db)
	createAccount := usecase3.NewCreateAccount(writer, authentication, session, user)
	account := handler3.NewAccount(createAccount)
	getOnboardingStep := usecase3.NewGetOnboardingStep(reader, user)
	userAttribute := repository2.NewUserAttribute(db)
	createUserAttribute := usecase3.NewCreateUserAttribute(writer, userAttribute)
	userProfile := repository2.NewUserProfile(db)
	createUserProfile := usecase3.NewCreateUserProfile(writer, userProfile)
	onboarding := handler3.NewOnboarding(getOnboardingStep, createUserAttribute, createUserProfile)
	usageCategory := repository3.NewUsageCategory(db)
	listUsageCategories := usecase3.NewListUsageCategories(reader, usageCategory)
	handlerUsageCategory := handler3.NewUsageCategory(listUsageCategories)
	handlers := di3.Handlers{
		Account:       account,
		Onboarding:    onboarding,
		UsageCategory: handlerUsageCategory,
	}
	return handlers
}

func InitUserRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di4.Repositories {
	user := repository2.NewUser(db)
	userAttribute := repository2.NewUserAttribute(db)
	userLinkProvider := repository2.NewUserLinkProvider(db)
	userLink := repository2.NewUserLink(db)
	userProfile := repository2.NewUserProfile(db)
	repositories := di4.Repositories{
		User:             user,
		UserAttribute:    userAttribute,
		UserLinkProvider: userLinkProvider,
		UserLink:         userLink,
		UserProfile:      userProfile,
	}
	return repositories
}

func InitUserUseCase(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di4.UseCases {
	userLink := repository2.NewUserLink(db)
	createUserLink := usecase4.NewCreateUserLink(writer, userLink)
	deleteUserLink := usecase4.NewDeleteUserLink(writer, userLink)
	listUserLink := usecase4.NewListUserLink(reader, userLink)
	updateUserLink := usecase4.NewUpdateUserLink(writer, userLink)
	useCases := di4.UseCases{
		CreateUserLink: createUserLink,
		DeleteUserLink: deleteUserLink,
		ListUserLink:   listUserLink,
		UpdateUserLink: updateUserLink,
	}
	return useCases
}

func InitUserHandler(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs *redis.Client) di4.Handlers {
	userLink := repository2.NewUserLink(db)
	createUserLink := usecase4.NewCreateUserLink(writer, userLink)
	listUserLink := usecase4.NewListUserLink(reader, userLink)
	updateUserLink := usecase4.NewUpdateUserLink(writer, userLink)
	deleteUserLink := usecase4.NewDeleteUserLink(writer, userLink)
	handlerUserLink := handler4.NewUserLink(createUserLink, listUserLink, updateUserLink, deleteUserLink)
	handlers := di4.Handlers{
		UserLink: handlerUserLink,
	}
	return handlers
}
