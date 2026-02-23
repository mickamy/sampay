package interceptor

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/mickamy/errx/cerr"

	ausecase "github.com/mickamy/sampay/internal/domain/auth/usecase"
	"github.com/mickamy/sampay/internal/misc/i18n"

	"github.com/mickamy/sampay/internal/di"
)

func NewInterceptors(infra *di.Infra) []connect.Interceptor {
	return []connect.Interceptor{
		Recovery(),
		cerr.NewInterceptor(
			cerr.WithLocaleFunc(func(header http.Header) string { return header.Get("Accept-Language") }),
			cerr.WithDefaultLocale(i18n.DefaultLanguage),
		),
		Logging(),
		I18N(),
		Authenticate(ausecase.NewAuthenticate(infra)),
	}
}
