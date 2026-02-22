package oauth

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/mickamy/sampay/internal/lib/random"
)

type Client interface {
	AuthenticationURL() (string, error)
	Callback(ctx context.Context, code string) (Payload, error)
}

func generateState(provider Provider) (string, error) {
	bytes, err := random.NewBytes(16)
	if err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	nonce := base64.URLEncoding.EncodeToString(bytes)
	return string(provider) + ":" + nonce, nil
}
