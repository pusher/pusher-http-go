package pusher

import (
	"encoding/json"
)

// Represents the information received about a channel from the Pusher API.
type Channel struct {
	Name              string
	Occupied          bool `json:"occupied,omitempty"`
	UserCount         int  `json:"user_count,omitempty"`
	SubscriptionCount int  `json:"subscription_count,omitempty"`
}

// Represents a list of channels received by the Pusher API.
type ChannelsList struct {
	Channels map[string]ChannelListItem `json:"channels"`
}

// An item of ChannelsList
type ChannelListItem struct {
	UserCount int `json:"user_count"`
}

// Represents a list of users in a presence-channel
type Users struct {
	List []User `json:"users"`
}

// Represents a user and contains their ID.
type User struct {
	Id string `json:"id"`
}

/*
A struct representing what to assign to a channel member, consisting of a `UserId` and any custom `UserInfo`.
*/
type MemberData struct {
	UserId   string            `json:"user_id"`
	UserInfo map[string]string `json:"user_info,omitempty"`
}

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
