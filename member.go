package pusher

import (
	"encoding/json"
)

type Member struct {
	UserId   string            `json:"user_id"`
	UserInfo map[string]string `json:"user_info,omitempty"`
}

func (m *Member) UserData() (userData string, err error) {
	var userDataBytes []byte
	if userDataBytes, err = json.Marshal(m); err != nil {
		return
	}
	userData = string(userDataBytes)
	return
}
