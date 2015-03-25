package pusher

import (
	"encoding/json"
	"errors"
)

type Event struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketId *string  `json:"socket_id,omitempty"`
}

type BufferedEvents struct {
	EventIds map[string]string `json:"event_ids,omitempty"`
}

func createTriggerPayload(channels []string, event string, data interface{}, socketId *string) ([]byte, error) {
	data2, _ := json.Marshal(data)

	dataAsString := string(data2)

	if len(dataAsString) > 10240 {
		return nil, errors.New("Data must be smaller than 10kb")
	}

	payload, _ := json.Marshal(&Event{
		Name:     event,
		Channels: channels,
		Data:     dataAsString,
		SocketId: socketId,
	})

	return payload, nil
}
