package cookie

import (
	"net/http"
	"time"

	"github.com/mickamy/sampay/config"
)

func Build(name, value string, expires time.Time) *http.Cookie {
	var maxAge int // 0 = session cookie (no Max-Age attribute)
	if !expires.IsZero() {
		remaining := time.Until(expires)
		if remaining > 0 {
			maxAge = int(remaining.Seconds())
		} else {
			maxAge = -1 // expired: delete cookie immediately
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
