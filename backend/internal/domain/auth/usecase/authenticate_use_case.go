package usecase

import (
	"context"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/auth/repository"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/lib/jwt"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var (
	ErrAuthenticateTokenNotSet = cmodel.NewLocalizableError(errx.NewSentinel("token not set", errx.Unauthenticated)).
					WithMessages(messages.AuthUseCaseErrorSessionNotSet())
	ErrAuthenticateTokenInvalid = cmodel.NewLocalizableError(errx.NewSentinel("token invalid", errx.Unauthenticated)).
					WithMessages(messages.AuthUseCaseErrorTokenInvalid())
	ErrAuthenticateSessionNotFound = cmodel.NewLocalizableError(
		errx.NewSentinel("session not found", errx.Unauthenticated)).
		WithMessages(messages.AuthUseCaseErrorSessionNotFound())
)

type AuthenticateInput struct {
	Token string
}

type AuthenticateOutput struct {
	UserID string
}

type Authenticate interface {
	Do(ctx context.Context, input AuthenticateInput) (AuthenticateOutput, error)
}

type authenticate struct {
	_           Authenticate       `inject:"returns"`
	_           *di.Infra          `inject:"param"`
	sessionRepo repository.Session `inject:""`
}

func (uc *authenticate) Do(ctx context.Context, input AuthenticateInput) (AuthenticateOutput, error) {
	if input.Token == "" {
		return AuthenticateOutput{}, ErrAuthenticateTokenNotSet
	}

	userID, err := jwt.ExtractID(input.Token)
	if err != nil {
		return AuthenticateOutput{}, ErrAuthenticateTokenInvalid
	}

	exists, err := uc.sessionRepo.AccessTokenExists(ctx, userID, input.Token)
	if err != nil {
		return AuthenticateOutput{}, errx.
			Wrap(err).
			With("message", "failed to check access token existence", "user_id", userID).
			WithCode(errx.Internal)
	}
	if !exists {
		return AuthenticateOutput{}, ErrAuthenticateSessionNotFound
	}

	return AuthenticateOutput{UserID: userID}, nil
}
