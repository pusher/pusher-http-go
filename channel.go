package pusher

import (
	"encoding/json"
)

type Channel struct {
	Name              string
	Client            Client
	Occupied          bool `json:"occupied"`
	UserCount         int  `json:"user_count",omitempty`
	SubscriptionCount int  `json:"subscription_count",omitempty`
}

func (c *Channel) GetUsers() (error, *Users) {
	path := "/apps/" + c.Client.AppId + "/channels/" + c.Name + "/users"
	u := Url{"GET", path, c.Client.Key, c.Client.Secret, nil, nil}
	err, raw_users := Request("GET", u.generate(), nil)

	users := &Users{}
	json.Unmarshal(raw_users, &users)
	return err, users
}
