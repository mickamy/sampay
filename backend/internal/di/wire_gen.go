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
	di2 "mickamy.com/sampay/internal/domain/common/di"
	handler2 "mickamy.com/sampay/internal/domain/common/handler"
	repository3 "mickamy.com/sampay/internal/domain/common/repository"
	usecase2 "mickamy.com/sampay/internal/domain/common/usecase"
	di3 "mickamy.com/sampay/internal/domain/registration/di"
	handler3 "mickamy.com/sampay/internal/domain/registration/handler"
	repository4 "mickamy.com/sampay/internal/domain/registration/repository"
	usecase3 "mickamy.com/sampay/internal/domain/registration/usecase"
	di4 "mickamy.com/sampay/internal/domain/user/di"
	handler4 "mickamy.com/sampay/internal/domain/user/handler"
	repository2 "mickamy.com/sampay/internal/domain/user/repository"
	usecase4 "mickamy.com/sampay/internal/domain/user/usecase"
	"mickamy.com/sampay/internal/job"
	"mickamy.com/sampay/internal/lib/aws/s3"
	"mickamy.com/sampay/internal/lib/aws/ses"
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
	infras := Infras{
		DB:         db,
		ReadWriter: readWriter,
		Writer:     writer,
		Reader:     reader,
		KVS:        v,
	}
	return infras, nil
}

func InitLibs() Libs {
	awsConfig := config.AWS()
	client := s3.New(awsConfig)
	sesClient := ses.New(awsConfig)
	libs := Libs{
		S3:  client,
		SES: sesClient,
	}
	return libs
}

func InitJobs() job.Jobs {
	awsConfig := config.AWS()
	client := ses.New(awsConfig)
	sendEmail := job.NewSendEmail(client)
	jobs := job.Jobs{
		SendEmail: sendEmail,
	}
	return jobs
}

func InitProducers() Producers {
	awsConfig := config.AWS()
	kvsConfig := config.KVS()
	producerConfig := provideProducerConfig(awsConfig, kvsConfig)
	client := provideSQSClient(awsConfig)
	producer := provideProducer(producerConfig, client)
	producers := Producers{
		Producer: producer,
	}
	return producers
}

func InitConsumers() Consumers {
	awsConfig := config.AWS()
	kvsConfig := config.KVS()
	consumerConfig := provideConsumerConfig(awsConfig, kvsConfig)
	jobs := InitJobs()
	consumer := provideConsumer(consumerConfig, jobs)
	consumers := Consumers{
		Consumer: consumer,
	}
	return consumers
}

func InitAuthRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di.Repositories {
	authentication := repository.NewAuthentication(db)
	emailVerification := repository.NewEmailVerification(db)
	session := repository.NewSession(kvs2)
	repositories := di.Repositories{
		Authentication:    authentication,
		EmailVerification: emailVerification,
		Session:           session,
	}
	return repositories
}

func InitAuthUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di.UseCases {
	emailVerification := repository.NewEmailVerification(db)
	authenticateAnonymousUser := usecase.NewAuthenticateAnonymousUser(reader, emailVerification)
	session := repository.NewSession(kvs2)
	user := repository2.NewUser(db)
	authenticateUser := usecase.NewAuthenticateUser(reader, session, user)
	authentication := repository.NewAuthentication(db)
	createSession := usecase.NewCreateSession(reader, authentication, session, user)
	deleteSession := usecase.NewDeleteSession(session)
	refreshSession := usecase.NewRefreshSession(session)
	awsConfig := config.AWS()
	kvsConfig := config.KVS()
	producerConfig := provideProducerConfig(awsConfig, kvsConfig)
	client := provideSQSClient(awsConfig)
	producer := provideProducer(producerConfig, client)
	requestEmailVerification := usecase.NewRequestEmailVerification(writer, producer, authentication, emailVerification)
	resetPassword := usecase.NewResetPassword(writer, emailVerification, authentication)
	verifyEmail := usecase.NewVerifyEmail(writer, producer, emailVerification, user, session)
	useCases := di.UseCases{
		AuthenticateAnonymousUser: authenticateAnonymousUser,
		AuthenticateUser:          authenticateUser,
		CreateSession:             createSession,
		DeleteSession:             deleteSession,
		RefreshSession:            refreshSession,
		RequestEmailVerification:  requestEmailVerification,
		ResetPassword:             resetPassword,
		VerifyEmail:               verifyEmail,
	}
	return useCases
}

func InitAuthHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di.Handlers {
	awsConfig := config.AWS()
	kvsConfig := config.KVS()
	producerConfig := provideProducerConfig(awsConfig, kvsConfig)
	client := provideSQSClient(awsConfig)
	producer := provideProducer(producerConfig, client)
	authentication := repository.NewAuthentication(db)
	emailVerification := repository.NewEmailVerification(db)
	requestEmailVerification := usecase.NewRequestEmailVerification(writer, producer, authentication, emailVerification)
	user := repository2.NewUser(db)
	session := repository.NewSession(kvs2)
	verifyEmail := usecase.NewVerifyEmail(writer, producer, emailVerification, user, session)
	handlerEmailVerification := handler.NewEmailVerification(requestEmailVerification, verifyEmail)
	resetPassword := usecase.NewResetPassword(writer, emailVerification, authentication)
	passwordReset := handler.NewPasswordReset(resetPassword)
	createSession := usecase.NewCreateSession(reader, authentication, session, user)
	refreshSession := usecase.NewRefreshSession(session)
	deleteSession := usecase.NewDeleteSession(session)
	handlerSession := handler.NewSession(createSession, refreshSession, deleteSession)
	handlers := di.Handlers{
		EmailVerification: handlerEmailVerification,
		PasswordReset:     passwordReset,
		Session:           handlerSession,
	}
	return handlers
}

func InitCommonRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di2.Repositories {
	s3Object := repository3.NewS3Object(db)
	repositories := di2.Repositories{
		S3Object: s3Object,
	}
	return repositories
}

func InitCommonUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di2.UseCases {
	awsConfig := config.AWS()
	client := s3.New(awsConfig)
	createDirectUploadURL := usecase2.NewCreateDirectUploadURL(client)
	useCases := di2.UseCases{
		CreateDirectUploadURL: createDirectUploadURL,
	}
	return useCases
}

func InitCommonHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di2.Handlers {
	awsConfig := config.AWS()
	client := s3.New(awsConfig)
	createDirectUploadURL := usecase2.NewCreateDirectUploadURL(client)
	directUploadURL := handler2.NewDirectUploadURL(createDirectUploadURL)
	handlers := di2.Handlers{
		DirectUploadURL: directUploadURL,
	}
	return handlers
}

func InitRegistrationRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di3.Repositories {
	usageCategory := repository4.NewUsageCategory(db)
	repositories := di3.Repositories{
		UsageCategory: usageCategory,
	}
	return repositories
}

func InitRegistrationUseCases(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di3.UseCases {
	authentication := repository.NewAuthentication(db)
	session := repository.NewSession(kvs2)
	user := repository2.NewUser(db)
	createAccount := usecase3.NewCreateAccount(writer, authentication, session, user)
	emailVerification := repository.NewEmailVerification(db)
	createPassword := usecase3.NewCreatePassword(writer, emailVerification, authentication)
	userAttribute := repository2.NewUserAttribute(db)
	createUserAttribute := usecase3.NewCreateUserAttribute(writer, userAttribute)
	userProfile := repository2.NewUserProfile(db)
	createUserProfile := usecase3.NewCreateUserProfile(writer, userProfile)
	getOnboardingStep := usecase3.NewGetOnboardingStep(reader, authentication, user)
	usageCategory := repository4.NewUsageCategory(db)
	listUsageCategories := usecase3.NewListUsageCategories(reader, usageCategory)
	useCases := di3.UseCases{
		CreateAccount:       createAccount,
		CreatePassword:      createPassword,
		CreateUserAttribute: createUserAttribute,
		CreateUserProfile:   createUserProfile,
		GetOnboardingStep:   getOnboardingStep,
		ListUsageCategories: listUsageCategories,
	}
	return useCases
}

func InitRegistrationHandlers(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di3.Handlers {
	authentication := repository.NewAuthentication(db)
	session := repository.NewSession(kvs2)
	user := repository2.NewUser(db)
	createAccount := usecase3.NewCreateAccount(writer, authentication, session, user)
	account := handler3.NewAccount(createAccount)
	getOnboardingStep := usecase3.NewGetOnboardingStep(reader, authentication, user)
	emailVerification := repository.NewEmailVerification(db)
	createPassword := usecase3.NewCreatePassword(writer, emailVerification, authentication)
	userAttribute := repository2.NewUserAttribute(db)
	createUserAttribute := usecase3.NewCreateUserAttribute(writer, userAttribute)
	userProfile := repository2.NewUserProfile(db)
	createUserProfile := usecase3.NewCreateUserProfile(writer, userProfile)
	onboarding := handler3.NewOnboarding(getOnboardingStep, createPassword, createUserAttribute, createUserProfile)
	usageCategory := repository4.NewUsageCategory(db)
	listUsageCategories := usecase3.NewListUsageCategories(reader, usageCategory)
	handlerUsageCategory := handler3.NewUsageCategory(listUsageCategories)
	handlers := di3.Handlers{
		Account:       account,
		Onboarding:    onboarding,
		UsageCategory: handlerUsageCategory,
	}
	return handlers
}

func InitUserRepositories(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di4.Repositories {
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

func InitUserUseCase(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di4.UseCases {
	userLink := repository2.NewUserLink(db)
	createUserLink := usecase4.NewCreateUserLink(writer, userLink)
	deleteUserLink := usecase4.NewDeleteUserLink(writer, userLink)
	user := repository2.NewUser(db)
	getMe := usecase4.NewGetMe(reader, user)
	getUser := usecase4.NewGetUser(reader, user)
	listUserLink := usecase4.NewListUserLink(reader, userLink)
	s3Object := repository3.NewS3Object(db)
	updateUserLinkQRCode := usecase4.NewUpdateUserLinkQRCode(writer, userLink, s3Object)
	updateUserLink := usecase4.NewUpdateUserLink(writer, userLink)
	userProfile := repository2.NewUserProfile(db)
	updateUserProfile := usecase4.NewUpdateUserProfile(writer, userProfile)
	updateUserProfileImage := usecase4.NewUpdateUserProfileImage(writer, userProfile, s3Object)
	useCases := di4.UseCases{
		CreateUserLink:         createUserLink,
		DeleteUserLink:         deleteUserLink,
		GetMe:                  getMe,
		GetUser:                getUser,
		ListUserLink:           listUserLink,
		UpdateUserLinkQRCode:   updateUserLinkQRCode,
		UpdateUserLink:         updateUserLink,
		UpdateUserProfile:      updateUserProfile,
		UpdateUserProfileImage: updateUserProfileImage,
	}
	return useCases
}

func InitUserHandler(db *database.DB, readWriter *database.ReadWriter, writer *database.Writer, reader *database.Reader, kvs2 *kvs.KVS) di4.Handlers {
	user := repository2.NewUser(db)
	getMe := usecase4.NewGetMe(reader, user)
	getUser := usecase4.NewGetUser(reader, user)
	handlerUser := handler4.NewUser(getMe, getUser)
	userLink := repository2.NewUserLink(db)
	createUserLink := usecase4.NewCreateUserLink(writer, userLink)
	listUserLink := usecase4.NewListUserLink(reader, userLink)
	updateUserLink := usecase4.NewUpdateUserLink(writer, userLink)
	s3Object := repository3.NewS3Object(db)
	updateUserLinkQRCode := usecase4.NewUpdateUserLinkQRCode(writer, userLink, s3Object)
	deleteUserLink := usecase4.NewDeleteUserLink(writer, userLink)
	handlerUserLink := handler4.NewUserLink(createUserLink, listUserLink, updateUserLink, updateUserLinkQRCode, deleteUserLink)
	userProfile := repository2.NewUserProfile(db)
	updateUserProfile := usecase4.NewUpdateUserProfile(writer, userProfile)
	updateUserProfileImage := usecase4.NewUpdateUserProfileImage(writer, userProfile, s3Object)
	handlerUserProfile := handler4.NewUserProfile(updateUserProfile, updateUserProfileImage)
	handlers := di4.Handlers{
		User:        handlerUser,
		UserLink:    handlerUserLink,
		UserProfile: handlerUserProfile,
	}
	return handlers
}
