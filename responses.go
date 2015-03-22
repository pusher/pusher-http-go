package pusher

import (
	"encoding/json"
)

func unmarshalledChannelsList(response []byte) *ChannelsList {
	channels := &ChannelsList{}
	json.Unmarshal(response, &channels)
	return channels
}

func unmarshalledChannel(response []byte, name string) *Channel {
	channel := &Channel{Name: name}
	json.Unmarshal(response, &channel)
	return channel
}

func unmarshalledChannelUsers(response []byte) *Users {
	users := &Users{}
	json.Unmarshal(response, &users)
	return users
}
