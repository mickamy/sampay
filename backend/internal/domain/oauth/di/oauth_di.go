package di

import (
	"github.com/google/wire"

	"mickamy.com/sampay/internal/domain/oauth/handler"
	"mickamy.com/sampay/internal/domain/oauth/usecase"
)

type UseCases struct {
	usecase.OAuthSignIn
	usecase.OAuthCallback
}

//lint:ignore U1000 used by wire
var UseCaseSet = wire.NewSet(
	usecase.NewOAuthSignIn,
	usecase.NewOAuthCallback,
)

type Handlers struct {
	*handler.OAuth
}

//lint:ignore U1000 used by wire
var HandlerSet = wire.NewSet(
	handler.NewOAuth,
)
