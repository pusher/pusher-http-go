package pusher

import (
	"encoding/json"
)

// This is a parsed form of a valid webhook received by the server.
type Webhook struct {
	// the timestamp of the request
	TimeMs int `json:"time_ms"`
	// the events associated with the webhook
	Events []WebhookEvent `json:"events"`
}

type WebhookEvent struct {
	// the type of the event
	Name string `json:"name"`
	// the channel on which it was sent
	Channel string `json:"channel"`
	// the name of the event
	Event string `json:"event,omitempty"`
	// the data associated with the event
	Data string `json:"data,omitempty"`
	// the socket_id of the sending socket
	SocketId string `json:"socket_id,omitempty"`
	// the user_id of a member who has joined or vacated a presence-channel
	UserId string `json:"user_id,omitempty"`
}

func newWebhook(body []byte) (*Webhook, error) {
	var webhook *Webhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		return nil, err
	}

	return webhook, nil
}
