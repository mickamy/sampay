package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"golang.org/x/text/language"

	"github.com/mickamy/sampay/internal/lib/logger"
	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/misc/i18n"
)

func I18N() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			header := req.Header().Get("Accept-Language")

			tags, _, err := language.ParseAcceptLanguage(header)
			if err != nil || len(tags) == 0 {
				logger.Warn(ctx, "failed to parse accept language tag. using default", "err", err, "Accept-Language", header)
				ctx = contexts.SetLanguage(ctx, i18n.DefaultLanguage)
				return next(ctx, req)
			}

			lang := i18n.ResolveLanguage(tags)
			ctx = contexts.SetLanguage(ctx, lang)
			return next(ctx, req)
		}
	}
}
