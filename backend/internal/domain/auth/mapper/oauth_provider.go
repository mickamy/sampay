package mapper

import (
	"github.com/mickamy/errx"

	authv1 "github.com/mickamy/sampay/gen/auth/v1"
	"github.com/mickamy/sampay/internal/domain/auth/model"
	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

func ToOAuthProvider(src authv1.OAuthProvider) (model.OAuthProvider, error) {
	switch src {
	case authv1.OAuthProvider_O_AUTH_PROVIDER_UNSPECIFIED:
		// fall through to error
	case authv1.OAuthProvider_O_AUTH_PROVIDER_GOOGLE:
		return model.OAuthProviderGoogle, nil
	case authv1.OAuthProvider_O_AUTH_PROVIDER_LINE:
		return model.OAuthProviderLINE, nil
	}
	return "", cmodel.NewLocalizableError(
		errx.New("unknown oauth provider", "provider", src).WithCode(errx.InvalidArgument),
	).
		WithMessages(messages.AuthMapperErrorUnknownOauthProvider())
}
