package messaging

import (
	"net/http"

	"github.com/mickamy/sampay/internal/domain/messaging/handler"
)

func RegisterWebhook(mux *http.ServeMux) {
	mux.Handle("/webhook/line", handler.NewWebhook())
}
