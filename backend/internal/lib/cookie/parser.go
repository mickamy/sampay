package cookie

import (
	"errors"
	"strings"
)

var (
	ErrNoAccessToken  = errors.New("cookie: no access token")
	ErrNoRefreshToken = errors.New("cookie: no refresh token")
)

func ParseAccessToken(cookie string) (string, error) {
	token := parse(cookie, "access_token")
	if token == "" {
		return "", ErrNoAccessToken
	}

	return token, nil
}

func ParseRefreshToken(cookie string) (string, error) {
	token := parse(cookie, "refresh_token")
	if token == "" {
		return "", ErrNoRefreshToken
	}

	return token, nil
}

func parse(cookie string, key string) string {
	parts := strings.Split(cookie, ";")
	cookieKey := key + "="
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if v, ok := strings.CutPrefix(part, cookieKey); ok {
			return v
		}
	}

	return ""
}
