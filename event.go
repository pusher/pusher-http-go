package pusher

import (
	"encoding/json"
	"errors"
)

type Event struct {
	Channel  string  `json:"channel"`
	Name     string  `json:"name"`
	Data     string  `json:"data"`
	SocketId *string `json:"socket_id,omitempty"`
}

type eventPayload struct {
	Name     string   `json:"name"`
	Channels []string `json:"channels"`
	Data     string   `json:"data"`
	SocketId *string  `json:"socket_id,omitempty"`
}

type BufferedEvents struct {
	EventIds map[string]string `json:"event_ids,omitempty"`
}

func createTriggerPayload(channels []string, event string, data interface{}, socketID *string, encryptionKey string) ([]byte, error) {
	var dataBytes []byte
	var err error
	var payloadData string

	switch d := data.(type) {
	case []byte:
		dataBytes = d
	case string:
		dataBytes = []byte(d)
	default:
		dataBytes, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}
	if isEncryptedChannel(channels[0]) {
		payloadData = encrypt(channels[0], dataBytes, encryptionKey)
	} else {
		payloadData = string(dataBytes)
	}

	if len(payloadData) > 10240 {
		return nil, errors.New("Data must be smaller than 10kb")
	}
	return json.Marshal(&eventPayload{
		Name:     event,
		Channels: channels,
		Data:     payloadData,
		SocketId: socketID,
	})
}
