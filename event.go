package pusher

import (
	"encoding/json"
	"errors"
)

type eventData struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketId *string  `json:"socket_id,omitempty"`
}

type BufferedEvents struct {
	EventIds map[string]string `json:"event_ids,omitempty"`
}

func createTriggerPayload(channels []string, event string, data interface{}, socketId *string) ([]byte, error) {
	data2, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if len(data2) > 10240 {
		return nil, errors.New("Data must be smaller than 10kb")
	}

	return json.Marshal(&Event{
		Name:     event,
		Channels: channels,
		Data:     string(data2),
		SocketId: socketId,
	})
}
