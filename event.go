package pusher

import (
	"encoding/json"
	"errors"
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

func createTriggerPayload(channels []string, event string, _data interface{}, socket_id string) ([]byte, error) {
	data, _ := json.Marshal(_data)

	data_as_string := string(data)

	if len(data_as_string) > 10240 {
		return nil, errors.New("Data must be smaller than 10kb")
	}

	payload, _ := json.Marshal(&Event{
		Name:     event,
		Channels: channels,
		Data:     data_as_string,
		SocketId: socket_id})

	return payload, nil
}
