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
