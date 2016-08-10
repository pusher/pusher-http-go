package pusher

import (
	"encoding/json"
)

// This is a parsed form of a valid webhook received by the server.
type Webhook struct {
	TimeMs int            `json:"time_ms"` // the timestamp of the request
	Events []WebhookEvent `json:"events"`  // the events associated with the webhook
}

type WebhookEvent struct {
	Name     string `json:"name"`                // the type of the event
	Channel  string `json:"channel"`             // the channel on which it was sent
	Event    string `json:"event,omitempty"`     // the name of the event
	Data     string `json:"data,omitempty"`      // the data associated with the event
	SocketId string `json:"socket_id,omitempty"` // the socket_id of the sending socket
	UserId   string `json:"user_id,omitempty"`   // the user_id of a member who has joined or vacated a presence-channel
}

func newWebhook(body []byte) (webhook *Webhook, err error) {
	err = json.Unmarshal(body, &webhook)
	return
}
