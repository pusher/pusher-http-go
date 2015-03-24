package pusher

import (
	"encoding/json"
)

type Webhook struct {
	TimeMs int            `json:"time_ms"`
	Events []WebhookEvent `json:"events"`
}

type WebhookEvent struct {
	Name     string `json:"name"`
	Channel  string `json:"channel"`
	Event    string `json:"event"`
	Data     string `json:"data"`
	SocketId string `json:"socket_id"`
}

func unmarshalledWebhook(request_body []byte) *Webhook {
	webhook := &Webhook{}
	json.Unmarshal(request_body, &webhook)
	return webhook
}
