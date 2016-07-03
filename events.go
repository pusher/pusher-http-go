package pusher

import (
	"encoding/json"
	"github.com/pusher/pusher/errors"
)

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

const maxDataSize = 10240

func (e *event) MarshalJSON() (body []byte, err error) {
	var dataJSON []byte

	switch d := e.Data.(type) {
	case []byte:
		dataJSON = d
	case string:
		dataJSON = []byte(d)
	default:
		if dataJSON, err = json.Marshal(d); err != nil {
			return
		}
	}

	if len(dataJSON) > maxDataSize {
		err = errors.New("Data must be smaller than 10kb")
		return
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
