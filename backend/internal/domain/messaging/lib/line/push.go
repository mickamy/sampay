package line

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mickamy/sampay/config"
)

const pushMessageURL = "https://api.line.me/v2/bot/message/push"

type PushClient struct {
	token string
}

func NewPushClient(cfg config.LineMessagingConfig) *PushClient {
	return &PushClient{token: cfg.ChannelToken}
}

type pushRequest struct {
	To       string        `json:"to"`
	Messages []pushMessage `json:"messages"`
}

type pushMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func (c *PushClient) SendTextMessage(ctx context.Context, to string, text string) error {
	body := pushRequest{
		To: to,
		Messages: []pushMessage{
			{Type: "text", Text: text},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("line push: failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, pushMessageURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("line push: failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := http.DefaultClient.Do(req) //nolint:gosec // URL is a constant, not user-controlled
	if err != nil {
		return fmt.Errorf("line push: failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("line push: unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
