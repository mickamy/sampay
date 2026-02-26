package i18n

import (
	"context"
	"embed"

	i18n "github.com/mickamy/go-typesafe-i18n"
	"golang.org/x/text/language"

	"github.com/mickamy/sampay/internal/misc/contexts"
)

//go:generate go tool go-typesafe-i18n generate -base=ja -out=./messages/messages_gen.go ./locales

//go:embed locales/*.yaml
var localeFS embed.FS

var bundle *i18n.Bundle
var DefaultLanguage = language.Japanese

func init() {
	bundle = i18n.NewBundle(DefaultLanguage)

	ja, _ := localeFS.ReadFile("locales/ja.yaml")
	bundle.MustLoadBytes("ja.yaml", ja)

	en, _ := localeFS.ReadFile("locales/en.yaml")
	bundle.MustLoadBytes("en.yaml", en)
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
