package contexts

import (
	"context"
	"errors"

	"mickamy.com/sampay/internal/domain/user/model"
	"mickamy.com/sampay/internal/lib/language"
)

type authenticatedUserKey struct{}
type languageKey struct{}

func SetAuthenticatedUser(ctx context.Context, user model.User) context.Context {
	return context.WithValue(ctx, authenticatedUserKey{}, user)
}

func AuthenticatedUser(ctx context.Context) (model.User, error) {
	user, ok := ctx.Value(authenticatedUserKey{}).(model.User)
	if ok {
		return user, nil
	}
	return user, errors.New("no authenticated user found in context")
}

func MustAuthenticatedUser(ctx context.Context) model.User {
	user, err := AuthenticatedUser(ctx)
	if err != nil {
		panic(err)
	}
	return user
}

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
