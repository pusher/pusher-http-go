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

type BufferedEvents struct {
	EventIds map[string]string `json:"event_ids"`
}

func createTriggerPayload(channels []string, event string, _data interface{}, socket_id string) ([]byte, error) {
	data, _ := json.Marshal(_data)

	dataAsString := string(data)

	if len(dataAsString) > 10240 {
		return nil, errors.New("Data must be smaller than 10kb")
	}

	payload, _ := json.Marshal(&Event{
		Name:     event,
		Channels: channels,
		Data:     dataAsString,
		SocketId: socket_id})

	return payload, nil
}
