package cookie

import (
	"net/http"
	"time"

	"github.com/mickamy/sampay/config"
)

func Build(name, value string, expires time.Time) *http.Cookie {
	maxAge := -1
	if expires.IsZero() {
		maxAge = -1
	} else {
		remaining := time.Until(expires)
		if remaining > 0 {
			maxAge = int(remaining.Seconds())
		}
	}
	var secure bool
	switch config.Common().Env {
	case config.EnvDevelopment, config.EnvTest:
		secure = false
	case config.EnvStaging, config.EnvProduction:
		secure = true
	}
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Secure:   secure,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
}
