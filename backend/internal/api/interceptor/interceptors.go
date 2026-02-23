package interceptor

import (
	"net/http"

	"connectrpc.com/connect"
	"github.com/mickamy/errx/cerr"
	"golang.org/x/text/language"

	ausecase "github.com/mickamy/sampay/internal/domain/auth/usecase"

	"github.com/mickamy/sampay/internal/di"
)

func NewInterceptors(infra *di.Infra) []connect.Interceptor {
	return []connect.Interceptor{
		Recovery(),
		cerr.NewInterceptor(
			cerr.WithLocaleFunc(func(header http.Header) string { return header.Get("Accept-Language") }),
			cerr.WithDefaultLocale(language.Japanese),
		),
		Logging(),
		I18N(),
		Authenticate(ausecase.NewAuthenticate(infra)),
	}
}
