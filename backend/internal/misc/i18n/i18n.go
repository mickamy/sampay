package i18n

import (
	"context"
	"path/filepath"

	i18n "github.com/mickamy/go-typesafe-i18n"
	"golang.org/x/text/language"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/misc/contexts"
)

//go:generate go tool go-typesafe-i18n generate -base=ja -out=./messages/messages_gen.go ./locales

var bundle *i18n.Bundle
var DefaultLanguage = language.Japanese

func init() {
	bundle = i18n.NewBundle(DefaultLanguage)
	locales := filepath.Join(config.Common().ModuleRoot, "internal", "misc", "i18n", "locales")
	bundle.MustLoadFile(filepath.Join(locales, "ja.yaml"))
	bundle.MustLoadFile(filepath.Join(locales, "en.yaml"))
}

func Localize(tag language.Tag, msg i18n.Message) string {
	return bundle.Localizer(tag).Localize(msg)
}

func LocalizeContext(ctx context.Context, msg i18n.Message) string {
	tag := contexts.MustLanguage(ctx)
	return Localize(tag, msg)
}

func Japanese(msg i18n.Message) string {
	return Localize(language.Japanese, msg)
}

var supportedLanguages = []language.Tag{
	language.Japanese,
	language.English,
}

var matcher = language.NewMatcher(supportedLanguages)

func ResolveLanguage(tags []language.Tag) language.Tag {
	if len(tags) == 0 {
		return DefaultLanguage
	}
	matched, _, confidence := matcher.Match(tags...)
	if confidence < language.High {
		return DefaultLanguage
	}
	return matched
}

type Message = i18n.Message
