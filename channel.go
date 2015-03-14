package pusher

type Channel struct {
	Name              string
	Occupied          bool `json:"occupied"`
	UserCount         int  `json:"user_count",omitempty`
	SubscriptionCount int  `json:"subscription_count",omitempty`
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
