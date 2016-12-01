package pusher

import (
	"encoding/json"
)

type Member struct {
	UserId   string            `json:"user_id"`
	UserInfo map[string]string `json:"user_info,omitempty"`
}

func (m *Member) UserData() (string, error) {
	userDataBytes, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(userDataBytes), nil
}
