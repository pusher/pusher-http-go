package pusher

import (
	"encoding/json"
	"github.com/pusher/pusher-http-go/errors"
)

// Max body size = 10 KB
const maxDataSize = 10240

type TriggerResponse struct {
	EventIds map[string]string `json:"event_ids,omitempty"`
}

type Event struct {
	Name     string `json:"name"`
	Channel  string `json:"channel"`
	Data     string `json:"data"`
	SocketID string `json:"socket_id,omitempty"`
}

type event struct {
	Name     string      `json:"name"`
	Channels []string    `json:"-"`
	Data     interface{} `json:"data"`
	SocketID *string     `json:"socket_id,omitempty"`
}

type batchRequest struct {
	Batch []Event `json:"batch"`
}

func (e *event) MarshalJSON() ([]byte, error) {
	var dataJSON []byte

	switch d := e.Data.(type) {
	case []byte:
		dataJSON = d
	case string:
		dataJSON = []byte(d)
	default:
		marshalled, err := json.Marshal(d)
		if err != nil {
			return nil, err
		}

		dataJSON = marshalled
	}

	if len(dataJSON) > maxDataSize {
		return nil, errors.New("Data must be smaller than 10kb")
	}

	return json.Marshal(&struct {
		Name     string   `json:"name"`
		Channels []string `json:"channels"`
		Data     string   `json:"data"`
		SocketID *string  `json:"socket_id,omitempty"`
	}{
		Name:     e.Name,
		Channels: e.Channels,
		Data:     string(dataJSON),
		SocketID: e.SocketID,
	})
}
