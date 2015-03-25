package pusher

import (
	"encoding/json"
)

type Channel struct {
	Name              string
	Occupied          bool `json:"occupied",omitempty`
	UserCount         int  `json:"user_count",omitempty`
	SubscriptionCount int  `json:"subscription_count",omitempty`
}

type ChannelsList struct {
	Channels map[string]ChannelListItem `json:"channels"`
}

type ChannelListItem struct {
	UserCount int `json:"user_count"`
}

type Users struct {
	List []User `json:"users"`
}

type User struct {
	Id string `json:"id"`
}

type MemberData struct {
	UserId   string            `json:"user_id"`
	UserInfo map[string]string `json:"user_info",omitempty`
}

// make sure to pass errors if any from json.Unmarshal

func unmarshalledChannelsList(response []byte) (*ChannelsList, error) {
	channels := &ChannelsList{}
	err := json.Unmarshal(response, &channels)

	if err != nil {
		return nil, err
	}

	return channels, nil
}

func unmarshalledChannel(response []byte, name string) (*Channel, error) {
	channel := &Channel{Name: name}
	err := json.Unmarshal(response, &channel)

	if err != nil {
		return nil, err
	}

	return channel, nil
}

func unmarshalledChannelUsers(response []byte) (*Users, error) {
	users := &Users{}
	err := json.Unmarshal(response, &users)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func unmarshalledBufferedEvents(response []byte) (*BufferedEvents, error) {
	bufferedEvents := &BufferedEvents{}
	err := json.Unmarshal(response, &bufferedEvents)

	if err != nil {
		return nil, err
	}

	return bufferedEvents, nil
}
