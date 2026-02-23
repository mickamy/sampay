package usecase

import (
	"context"

	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/domain/auth/repository"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/lib/jwt"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var (
	ErrRefreshTokenNotSet = cmodel.NewLocalizableError(errx.NewSentinel("token not set", errx.InvalidArgument)).
				WithMessages(messages.AuthUseCaseErrorSessionNotSet())
	ErrRefreshTokenInvalid = cmodel.NewLocalizableError(errx.NewSentinel("token invalid", errx.InvalidArgument)).
				WithMessages(messages.AuthUseCaseErrorTokenInvalid())
	ErrRefreshTokenNotFound = cmodel.NewLocalizableError(errx.NewSentinel("session not found", errx.InvalidArgument)).
				WithMessages(messages.AuthUseCaseErrorSessionNotFound())
)

type RefreshTokenInput struct {
	Token string
}

type RefreshTokenOutput struct {
	Tokens jwt.Tokens
}

type RefreshToken interface {
	Do(ctx context.Context, input RefreshTokenInput) (RefreshTokenOutput, error)
}

type refreshToken struct {
	_           RefreshToken       `inject:"returns"`
	_           *di.Infra          `inject:"param"`
	sessionRepo repository.Session `inject:""`
}

func (uc *refreshToken) Do(ctx context.Context, input RefreshTokenInput) (RefreshTokenOutput, error) {
	if input.Token == "" {
		return RefreshTokenOutput{}, ErrRefreshTokenNotSet
	}

	userID, err := jwt.ExtractID(input.Token)
	if err != nil {
		return RefreshTokenOutput{}, ErrRefreshTokenInvalid
	}

	exists, err := uc.sessionRepo.RefreshTokenExists(ctx, userID, input.Token)
	if err != nil {
		return RefreshTokenOutput{}, errx.Wrap(err, "failed to check session existence")
	}
	if !exists {
		return RefreshTokenOutput{}, ErrRefreshTokenNotFound
	}

	session, err := model.NewSession(userID)
	if err != nil {
		return RefreshTokenOutput{}, errx.Wrap(err, "failed to initialize session")
	}

	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return RefreshTokenOutput{}, errx.Wrap(err, "failed to create session")
	}

	return RefreshTokenOutput{Tokens: session.Tokens}, nil
}
