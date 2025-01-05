package i18n

import (
	"context"
	_ "embed"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/yaml.v3"

	"mickamy.com/sampay/internal/lib/contexts"
	"mickamy.com/sampay/internal/lib/language"
)

var bundle *i18n.Bundle

//go:embed resources/ja.yaml
var japanese []byte

func init() {
	bundle = i18n.NewBundle(language.Japanese.Tag())
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	bundle.MustParseMessageFileBytes(japanese, "ja.yaml")
}

type Config struct {
	// MessageID is the id of the message to lookup.
	// This field is ignored if DefaultMessage is set.
	MessageID string

	// TemplateData is the data passed when executing the message's template.
	// If TemplateData is nil and PluralCount is not nil, then the message template
	// will be executed with data that contains the plural count.
	TemplateData interface{}

	// PluralCount determines which plural form of the message is used.
	PluralCount interface{}
}

func (cfg Config) convert() *i18n.LocalizeConfig {
	return &i18n.LocalizeConfig{
		MessageID:    cfg.MessageID,
		TemplateData: cfg.TemplateData,
		PluralCount:  cfg.PluralCount,
	}
}

func MustLocalizeMessage(lang language.Type, config Config) string {
	localizer := i18n.NewLocalizer(bundle, string(lang))
	message := localizer.MustLocalize(config.convert())
	return message
}

func MustLocalizeMessageCtx(ctx context.Context, config Config) string {
	lang := contexts.MustLanguage(ctx)
	return MustLocalizeMessage(lang, config)
}
