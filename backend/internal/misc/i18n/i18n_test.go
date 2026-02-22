package i18n_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/mickamy/sampay/internal/misc/contexts"
	"github.com/mickamy/sampay/internal/misc/i18n"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

//nolint:gosmopolitan // intentional i18n test
func TestLocalize(t *testing.T) {
	t.Parallel()

	jaMsg := i18n.Localize(language.Japanese, messages.CommonResponseErrorInternal())
	assert.Equal(t,
		"ただいまアクセスが集中しております。しばらくしてから再度お試しください。", jaMsg)
	enMsg := i18n.Localize(language.English, messages.CommonResponseErrorInternal())
	assert.Equal(t,
		"We are currently experiencing high traffic. Please try again later.", enMsg)
}

//nolint:gosmopolitan // intentional i18n test
func TestLocalizeContext(t *testing.T) {
	t.Parallel()

	ctx := t.Context()

	{
		ctx = contexts.SetLanguage(ctx, language.Japanese)
		jaMsg := i18n.LocalizeContext(ctx, messages.CommonResponseErrorInternal())
		assert.Equal(t,
			"ただいまアクセスが集中しております。しばらくしてから再度お試しください。", jaMsg)
	}
	{
		ctx = contexts.SetLanguage(ctx, language.English)
		enMsg := i18n.LocalizeContext(ctx, messages.CommonResponseErrorInternal())
		assert.Equal(t,
			"We are currently experiencing high traffic. Please try again later.", enMsg)
	}
}

func TestResolveLanguage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []language.Tag
		want language.Tag
	}{
		{
			name: "Japanese when Japanese is first",
			in:   []language.Tag{language.Japanese, language.English},
			want: language.Japanese,
		},
		{
			name: "English when English is first",
			in:   []language.Tag{language.English, language.Japanese},
			want: language.English,
		},
		{
			name: "skip unsupported and return first supported",
			in:   []language.Tag{language.French, language.English},
			want: language.English,
		},
		{
			name: "default to Japanese when no supported language",
			in:   []language.Tag{language.French, language.German},
			want: language.Japanese,
		},
		{
			name: "default to Japanese when empty",
			in:   []language.Tag{},
			want: language.Japanese,
		},
		{
			name: "default to Japanese when nil",
			in:   nil,
			want: language.Japanese,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, i18n.ResolveLanguage(tt.in))
		})
	}
}
