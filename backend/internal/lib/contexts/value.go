package contexts

import (
	"context"
	"errors"

	"mickamy.com/sampay/internal/lib/language"
)

type authenticatedUserIDKey struct{}
type anonymousUserTokenKey struct{}
type languageKey struct{}

func SetAuthenticatedUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, authenticatedUserIDKey{}, userID)
}

func AuthenticatedUserID(ctx context.Context) (string, error) {
	id, ok := ctx.Value(authenticatedUserIDKey{}).(string)
	if ok {
		return id, nil
	}
	return id, errors.New("no authenticated user id found in context")
}

func MustAuthenticatedUserID(ctx context.Context) string {
	id, err := AuthenticatedUserID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

func SetAnonymousUserToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, anonymousUserTokenKey{}, token)
}

func AnonymousUserToken(ctx context.Context) (string, error) {
	token, ok := ctx.Value(anonymousUserTokenKey{}).(string)
	if ok {
		return token, nil
	}
	return token, errors.New("no anonymous user token found in context")
}

func MustAnonymousUserToken(ctx context.Context) string {
	token, err := AnonymousUserToken(ctx)
	if err != nil {
		panic(err)
	}
	return token
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
