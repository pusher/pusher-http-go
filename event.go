package pusher

import (
	"encoding/json"
	"errors"
)

// maxEventPayloadSize indicates the max size allowed for the data content (payload) of each event
const maxEventPayloadSize = 10240

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

	dataBytes, err := byteEncodePayload(data)
	if err != nil {
		return nil, err
	}

	var payloadData string
	if isEncryptedChannel(channels[0]) {
		payloadData = encrypt(channels[0], dataBytes, encryptionKey)
	} else {
		payloadData = string(dataBytes)
	}

	if len(payloadData) > maxEventPayloadSize {
		return nil, errors.New("Data must be smaller than 10kb")
	}
	return json.Marshal(&eventPayload{
		Name:     event,
		Channels: channels,
		Data:     payloadData,
		SocketId: socketID,
	})
}

func byteEncodePayload(data interface{}) ([]byte, error) {
	var dataBytes []byte
	var err error

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
	return dataBytes, nil
}
