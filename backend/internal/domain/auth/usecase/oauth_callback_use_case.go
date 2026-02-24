package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mickamy/errx"

	"github.com/mickamy/sampay/internal/di"
	"github.com/mickamy/sampay/internal/domain/auth/model"
	"github.com/mickamy/sampay/internal/domain/auth/repository"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	umodel "github.com/mickamy/sampay/internal/domain/user/model"
	urepository "github.com/mickamy/sampay/internal/domain/user/repository"
	"github.com/mickamy/sampay/internal/infra/storage/database"
	"github.com/mickamy/sampay/internal/lib/oauth"
	"github.com/mickamy/sampay/internal/lib/ulid"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

var (
	ErrOAuthCallbackUnsupportedProvider = cmodel.NewLocalizableError(
		errx.NewSentinel("unsupported provider", errx.InvalidArgument)).
		WithMessages(messages.AuthUseCaseErrorUnsupportedOauthProvider())
	ErrOAuthCallbackFailed = cmodel.NewLocalizableError(errx.NewSentinel("oauth callback failed", errx.InvalidArgument)).
				WithMessages(messages.AuthUseCaseErrorOauthCallbackFailed())
)

type OAuthCallbackInput struct {
	Provider model.OAuthProvider
	Code     string
}

type OAuthCallbackOutput struct {
	Session   model.Session
	EndUser   umodel.EndUser
	IsNewUser bool
}

type OAuthCallback interface {
	Do(ctx context.Context, input OAuthCallbackInput) (OAuthCallbackOutput, error)
}

type oauthCallback struct {
	_                OAuthCallback           `inject:"returns"`
	_                *di.Infra               `inject:"param"`
	resolver         *oauth.Resolver         `inject:"param"`
	writer           *database.Writer        `inject:""`
	userRepo         urepository.User        `inject:""`
	endUserRepo      urepository.EndUser     `inject:""`
	oauthAccountRepo repository.OAuthAccount `inject:""`
	sessionRepo      repository.Session      `inject:""`
}

func (uc *oauthCallback) Do(ctx context.Context, input OAuthCallbackInput) (OAuthCallbackOutput, error) {
	client, err := uc.resolveClient(input.Provider)
	if err != nil {
		return OAuthCallbackOutput{}, err
	}

	payload, err := client.Callback(ctx, input.Code)
	if err != nil {
		return OAuthCallbackOutput{}, errors.Join(ErrOAuthCallbackFailed, err)
	}

	var session model.Session
	var endUser umodel.EndUser
	var isNewUser bool
	if err := uc.writer.Transaction(ctx, func(tx *database.DB) error {
		existingAccount, err := uc.oauthAccountRepo.WithTx(tx).GetByProviderAndUID(ctx, input.Provider, payload.UID)
		if err != nil && !errors.Is(err, database.ErrNotFound) {
			return errx.Wrap(err, "message", "failed to get existing account").
				WithCode(errx.Internal)
		}

		if errors.Is(err, database.ErrNotFound) {
			userID := ulid.New()
			baseUser := umodel.User{ID: userID}
			if err := uc.userRepo.WithTx(tx).Create(ctx, &baseUser); err != nil {
				return errx.Wrap(err, "message", "failed to create user").
					WithCode(errx.Internal)
			}

			endUser = umodel.EndUser{
				UserID: userID,
				Slug:   uuid.NewString(),
			}
			if err := uc.endUserRepo.WithTx(tx).Create(ctx, &endUser); err != nil {
				return errx.Wrap(err, "message", "failed to create end user").
					WithCode(errx.Internal)
			}

			oauthAccount := model.OAuthAccount{
				ID:        ulid.New(),
				EndUserID: userID,
				Provider:  payload.Provider.String(),
				UID:       payload.UID,
			}
			if err := uc.oauthAccountRepo.WithTx(tx).Create(ctx, &oauthAccount); err != nil {
				return errx.Wrap(err, "message", "failed to create oauth account").
					WithCode(errx.Internal)
			}

			isNewUser = true
		} else {
			endUser, err = uc.endUserRepo.WithTx(tx).Get(ctx, existingAccount.EndUserID)
			if err != nil {
				return errx.Wrap(err, "message", "failed to get end user").
					WithCode(errx.Internal)
			}
		}

		session, err = model.NewSession(endUser.UserID)
		if err != nil {
			return errx.Wrap(err, "message", "failed to initialize session").
				WithCode(errx.Internal)
		}

		if err := uc.sessionRepo.Create(ctx, session); err != nil {
			return errx.Wrap(err, "message", "failed to create session").
				WithCode(errx.Internal)
		}

		return nil
	}); err != nil {
		return OAuthCallbackOutput{}, err //nolint:wrapcheck // errors from transaction callback are already wrapped inside
	}

	return OAuthCallbackOutput{Session: session, EndUser: endUser, IsNewUser: isNewUser}, nil
}

func (uc *oauthCallback) resolveClient(provider model.OAuthProvider) (oauth.Client, error) {
	client, err := uc.resolver.Resolve(oauth.Provider(provider))
	if err != nil {
		return nil, errx.Wrap(ErrOAuthCallbackUnsupportedProvider, "provider", provider, "err", err)
	}
	return client, nil
}
