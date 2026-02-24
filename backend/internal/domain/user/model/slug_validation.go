package model

import (
	"regexp"

	"github.com/mickamy/errx"

	cmodel "github.com/mickamy/sampay/internal/domain/common/model"
	"github.com/mickamy/sampay/internal/misc/i18n"
	"github.com/mickamy/sampay/internal/misc/i18n/messages"
)

const (
	slugMinLen = 3
	slugMaxLen = 30
)

var (
	slugPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)
	uuidPattern = regexp.MustCompile(
		`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`,
	)
)

var reservedSlugs = map[string]struct{}{
	"admin":    {},
	"api":      {},
	"auth":     {},
	"login":    {},
	"logout":   {},
	"my":       {},
	"oauth":    {},
	"setup":    {},
	"settings": {},
	"support":  {},
	"help":     {},
}

// ValidateSlug validates the given slug and returns localised error messages.
// Returns nil when the slug is valid.
func ValidateSlug(slug string) *cmodel.LocalizableError {
	var msgs []i18n.Message

	if len(slug) < slugMinLen {
		msgs = append(msgs, messages.UserModelSlugErrorTooShort())
	}
	if len(slug) > slugMaxLen {
		msgs = append(msgs, messages.UserModelSlugErrorTooLong())
	}
	if len(slug) >= slugMinLen && !slugPattern.MatchString(slug) {
		msgs = append(msgs, messages.UserModelSlugErrorInvalidChars())
	}
	if uuidPattern.MatchString(slug) {
		msgs = append(msgs, messages.UserModelSlugErrorIsUuid())
	}
	if _, ok := reservedSlugs[slug]; ok {
		msgs = append(msgs, messages.UserModelSlugErrorReserved())
	}

	if len(msgs) == 0 {
		return nil
	}

	return cmodel.NewLocalizableError(errSlugValidation).WithMessages(msgs...)
}

var errSlugValidation = cmodel.NewLocalizableError(errx.NewSentinel("slug validation failed", errx.InvalidArgument))
