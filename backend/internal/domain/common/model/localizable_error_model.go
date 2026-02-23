package model

import (
	"context"
	"strings"

	"github.com/mickamy/errx"
	"golang.org/x/text/language"

	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/misc/i18n"
)

type LocalizableError struct {
	underlying error
	messages   []i18n.Message
}

func NewLocalizableError(underlying error) *LocalizableError {
	return &LocalizableError{
		underlying: underlying,
	}
}

func (m *LocalizableError) Error() string { return m.underlying.Error() }

func (m *LocalizableError) Unwrap() error { return m.underlying }

func (m *LocalizableError) WithMessages(messages ...i18n.Message) *LocalizableError {
	m.messages = append(m.messages, messages...)
	return m
}

func (m *LocalizableError) Localize(locale string) string {
	lang, err := language.Parse(locale)
	if err != nil {
		lang = i18n.DefaultLanguage
	}
	return m.LocalizeTag(lang)
}

func (m *LocalizableError) LocalizeTag(lang language.Tag) string {
	var localized string
	for _, message := range m.messages {
		if localized == "" {
			localized = i18n.Localize(lang, message)
		} else {
			localized = strings.Join([]string{localized, i18n.Localize(lang, message)}, "\n")
		}
	}
	return localized
}

func (m *LocalizableError) LocalizeContext(ctx context.Context) string {
	lang := contexts.MustLanguage(ctx)
	return m.LocalizeTag(lang)
}

var _ error = (*LocalizableError)(nil)
var _ errx.Localizable = (*LocalizableError)(nil)
