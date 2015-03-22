package pusher

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
	UserId   int               `json:"user_id"`
	UserInfo map[string]string `json:"user_info",omitempty`
}
