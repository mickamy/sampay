package interceptor

import (
	"context"

	"connectrpc.com/connect"
	"github.com/mickamy/slogger"
	txtlang "golang.org/x/text/language"

	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/language"
)

func I18N() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			header := req.Header().Get("Accept-Language")

			tag, _, err := txtlang.ParseAcceptLanguage(header)
			if err != nil || len(tag) == 0 {
				slogger.WarnCtx(ctx, "failed to parse accept language tag", "err", err, "Accept-Language", header)
				ctx = contexts.SetLanguage(ctx, language.Japanese)
				return next(ctx, req)
			}

			lang := language.Type(tag[0].String())
			if !lang.IsSupported() {
				slogger.InfoCtx(ctx, "unsupported language", "lang", lang)
				lang = language.Japanese
			}

			ctx = contexts.SetLanguage(ctx, lang)
			return next(ctx, req)
		}
	}
}
