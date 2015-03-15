package pusher

import (
	"net/http"
)

type Webhook struct {
	TimeMs      int            `json:"time_ms"`
	Events      []WebhookEvent `json:"events"`
	Key, Secret string
	Header      http.Header
	RawBody     string
}

func (w *Webhook) IsValid() bool {

	for _, token := range w.Header["X-Pusher-Key"] {
		if token == w.Key {
			return CheckSignature(w.Header["X-Pusher-Signature"][0], w.RawBody, w.Secret)
		}
	}
	return false

}
