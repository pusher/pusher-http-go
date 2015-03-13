package pusher

import (
	"encoding/json"
	"fmt"
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
	q := Url{"GET", path, c.Client.Key, c.Client.Secret, nil, nil}
	err, raw_users := c.Client.get(q.generate(), nil)

	fmt.Println(string(raw_users))

	users := &Users{}
	json.Unmarshal(raw_users, &users)
	return err, users
}
