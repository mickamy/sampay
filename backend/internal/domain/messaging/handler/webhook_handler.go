package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/mickamy/sampay/config"
	"github.com/mickamy/sampay/internal/lib/logger"
)

type Webhook struct {
	channelSecret string
}

func NewWebhook() *Webhook {
	return &Webhook{
		channelSecret: config.LineMessaging().ChannelSecret,
	}
}

func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(r.Context(), "failed to read webhook body", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	signature := r.Header.Get("X-Line-Signature")
	if !h.verifySignature(body, signature) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	logger.Info(r.Context(), "received LINE webhook")

	w.WriteHeader(http.StatusOK)
}

func (h *Webhook) verifySignature(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(h.channelSecret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
