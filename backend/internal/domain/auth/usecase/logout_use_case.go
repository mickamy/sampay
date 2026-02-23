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
	ErrLogoutInvalidAccessToken = cmodel.NewLocalizableError(errx.NewSentinel("token invalid", errx.InvalidArgument)).
					WithMessages(messages.AuthUseCaseErrorTokenInvalid())
	ErrLogoutInvalidRefreshToken = cmodel.NewLocalizableError(errx.NewSentinel("token invalid", errx.InvalidArgument)).
					WithMessages(messages.AuthUseCaseErrorTokenInvalid())
	ErrLogoutTokenMismatch = cmodel.NewLocalizableError(errx.NewSentinel("token mismatch", errx.InvalidArgument)).
				WithMessages(messages.AuthUseCaseErrorLogoutTokenMismatch())
)

type LogoutInput struct {
	AccessToken  string //nolint:gosec // internal use only, not serialized to JSON
	RefreshToken string //nolint:gosec // internal use only, not serialized to JSON
}

type LogoutOutput struct{}

type Logout interface {
	Do(ctx context.Context, input LogoutInput) (LogoutOutput, error)
}

type logout struct {
	_           Logout             `inject:"returns"`
	_           *di.Infra          `inject:"param"`
	sessionRepo repository.Session `inject:""`
}

func (uc *logout) Do(ctx context.Context, input LogoutInput) (LogoutOutput, error) {
	userID, err := jwt.ExtractID(input.AccessToken)
	if err != nil {
		return LogoutOutput{}, errx.Wrap(
			ErrLogoutInvalidAccessToken,
			"failed to extract user id from access token",
			"err",
			err,
		)
	}

	userIDFromRefresh, err := jwt.ExtractID(input.RefreshToken)
	if err != nil {
		return LogoutOutput{}, errx.Wrap(
			ErrLogoutInvalidRefreshToken,
			"failed to extract user id from refresh token",
			"err",
			err,
		)
	}

	if userID != userIDFromRefresh {
		return LogoutOutput{}, errx.Wrap(
			ErrLogoutTokenMismatch,
			"access", userID,
			"refresh", userIDFromRefresh,
		)
	}

	if err := uc.sessionRepo.Delete(ctx, model.Session{
		UserID: userID,
		Tokens: jwt.Tokens{
			Access:  jwt.Token{Value: input.AccessToken},
			Refresh: jwt.Token{Value: input.RefreshToken},
		},
	}); err != nil {
		return LogoutOutput{}, errx.Wrap(err, "failed to delete session", "user_id", userID)
	}

	return LogoutOutput{}, nil
}
