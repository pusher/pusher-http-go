package pusher

// Represents the information received about a channel from the Pusher API.
type Channel struct {
	Occupied          bool `json:"occupied"`
	UserCount         int  `json:"user_count,omitempty"`
	SubscriptionCount int  `json:"subscription_count,omitempty"`
}

// Represents a list of channels received by the Pusher API.
type ChannelList struct {
	Channels map[string]ChannelListItem `json:"channels"`
}

// An item of ChannelsList
type ChannelListItem struct {
	UserCount int `json:"user_count,omitempty"`
}

type UserList struct {
	Users []User `json:"users"`
}

type User struct {
	Id string `json:"id"`
}
