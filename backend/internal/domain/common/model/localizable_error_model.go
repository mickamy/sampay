package model

import (
	"context"
	"strings"

	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/language"
	"mickamy.com/sampay/internal/misc/i18n"
)

type LocalizableError struct {
	UnderlyingError error

	Messages []i18n.Config
}

func NewLocalizableError(underlyingError error) *LocalizableError {
	return &LocalizableError{
		UnderlyingError: underlyingError,
	}
}

func (m *LocalizableError) WithMessages(messages ...i18n.Config) *LocalizableError {
	m.Messages = append(m.Messages, messages...)
	return m
}

func (m *LocalizableError) Localize(lang language.Type) string {
	var localized string
	for _, message := range m.Messages {
		if localized == "" {
			localized = i18n.MustLocalizeMessage(lang, message)
		} else {
			localized = strings.Join([]string{localized, i18n.MustLocalizeMessage(lang, message)}, "\n")
		}
	}
	return localized
}

func (m *LocalizableError) LocalizeCtx(ctx context.Context) string {
	lang := contexts.MustLanguage(ctx)
	return m.Localize(lang)
}

func (m *LocalizableError) Error() string {
	return m.UnderlyingError.Error()
}
