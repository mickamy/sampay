package contexts

import (
	"context"
	"errors"

	"mickamy.com/sampay/internal/lib/language"
)

type languageKey struct{}

func SetLanguage(ctx context.Context, lang language.Type) context.Context {
	return context.WithValue(ctx, languageKey{}, lang)
}

func Language(ctx context.Context) (language.Type, error) {
	lang, ok := ctx.Value(languageKey{}).(language.Type)
	if ok {
		return lang, nil
	}
	return lang, errors.New("no language found in context")
}

func MustLanguage(ctx context.Context) language.Type {
	lang, err := Language(ctx)
	if err != nil {
		panic(err)
	}
	return lang
}
