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

//return error

func unmarshalledWebhook(requestBody []byte) (*Webhook, error) {
	webhook := &Webhook{}
	err := json.Unmarshal(requestBody, &webhook)
	if err != nil {
		return nil, err
	}
	return webhook, nil
}
