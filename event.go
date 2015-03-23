package pusher

import (
	"encoding/json"
)

type Event struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketId string   `json:"socket_id"`
}

type WebhookEvent struct {
	Name     string `json:"name"`
	Channel  string `json:"channel"`
	Event    string `json:"event"`
	Data     string `json:"data"`
	SocketId string `json:"socket_id"`
}

func createTriggerPayload(channels []string, event string, _data interface{}, socket_id string) []byte {
	data, _ := json.Marshal(_data)

	payload, _ := json.Marshal(&Event{
		Name:     event,
		Channels: channels,
		Data:     string(data),
		SocketId: socket_id})

	return payload
}
