package pusher

import (
	"encoding/json"
)

// Channel represents the information about a channel from the Pusher API.
type Channel struct {
	Name              string
	Occupied          bool `json:"occupied,omitempty"`
	UserCount         int  `json:"user_count,omitempty"`
	SubscriptionCount int  `json:"subscription_count,omitempty"`
}

// ChannelsList represents a list of channels received by the Pusher API.
type ChannelsList struct {
	Channels map[string]ChannelListItem `json:"channels"`
}

// ChannelListItem represents an item within ChannelsList
type ChannelListItem struct {
	UserCount int `json:"user_count"`
}

// Users represents a list of users in a presence-channel
type Users struct {
	List []User `json:"users"`
}

// User represents a user and contains their ID.
type User struct {
	ID string `json:"id"`
}

/*
MemberData represents what to assign to a channel member, consisting of a
`UserID` and any custom `UserInfo`.
*/
type MemberData struct {
	UserID   string            `json:"user_id"`
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
