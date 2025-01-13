package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/common/handler"
	"mickamy.com/sampay/internal/domain/common/repository"
	"mickamy.com/sampay/internal/domain/common/usecase"
)

type Repositories struct {
	repository.S3Object
}

//lint:ignore U1000 used by wire
var RepositorySet = wire.NewSet(
	repository.NewS3Object,
)

type UseCases struct {
	usecase.CreateDirectUploadURL
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewCreateDirectUploadURL,
)

type Handlers struct {
	*handler.DirectUploadURL
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewDirectUploadURL,
)
