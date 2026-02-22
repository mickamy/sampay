package contexts

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/text/language"
)

type authenticatedUserIDKey struct{}
type executionIDKey struct{}
type systemUserIDKey struct{}
type languageKey struct{}

// SetAuthenticatedUserID sets the authenticated user ID in the context.
func SetAuthenticatedUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, authenticatedUserIDKey{}, userID)
}

// AuthenticatedUserID retrieves the authenticated user ID from the context.
func AuthenticatedUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(authenticatedUserIDKey{}).(string)
	if ok {
		return userID, nil
	}
	return "", errors.New("contexts: no authenticated user ID found in context")
}

// MustAuthenticatedUserID retrieves the authenticated user ID from the context and panics if not found.
func MustAuthenticatedUserID(ctx context.Context) string {
	userID, err := AuthenticatedUserID(ctx)
	if err != nil {
		panic(err)
	}
	return userID
}

// SetExecutionID sets the execution ID in the context.
func SetExecutionID(ctx context.Context, requestID uuid.UUID) context.Context {
	return context.WithValue(ctx, executionIDKey{}, requestID)
}

// ExecutionID retrieves the execution ID from the context.
func ExecutionID(ctx context.Context) (uuid.UUID, error) {
	requestID, ok := ctx.Value(executionIDKey{}).(uuid.UUID)
	if ok {
		return requestID, nil
	}
	return uuid.Nil, errors.New("contexts: no execution ID found in context")
}

// MustExecutionID retrieves the execution ID from the context and panics if not found.
func MustExecutionID(ctx context.Context) uuid.UUID {
	requestID, err := ExecutionID(ctx)
	if err != nil {
		panic(err)
	}
	return requestID
}

// SetSystemUserID sets the system user ID in the context.
func SetSystemUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, systemUserIDKey{}, userID)
}

// SystemUserID retrieves the system user ID from the context.
func SystemUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(systemUserIDKey{}).(string)
	if ok {
		return userID, nil
	}
	return "", errors.New("contexts: no system user ID found in context")
}

// MustSystemUserID retrieves the system user ID from the context and panics if not found.
func MustSystemUserID(ctx context.Context) string {
	userID, err := SystemUserID(ctx)
	if err != nil {
		panic(err)
	}
	return userID
}

// SetLanguage sets the language in the context.
func SetLanguage(ctx context.Context, lang language.Tag) context.Context {
	return context.WithValue(ctx, languageKey{}, lang)
}

// Language retrieves the language from the context.
func Language(ctx context.Context) (language.Tag, error) {
	lang, ok := ctx.Value(languageKey{}).(language.Tag)
	if ok {
		return lang, nil
	}
	return lang, errors.New("contexts: no language found in context")
}

// MustLanguage retrieves the language from the context and panics if not found.
func MustLanguage(ctx context.Context) language.Tag {
	lang, err := Language(ctx)
	if err != nil {
		panic(err)
	}
	return lang
}
