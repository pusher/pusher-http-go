package authentications

import (
	"fmt"
)

type PresenceChannel struct {
	Body []byte
	Member
}

func (p *PresenceChannel) StringToSign() (string, error) {
	privateChannelRequest := &PrivateChannel{p.Body}
	unsigned, err := privateChannelRequest.StringToSign()
	if err != nil {
		return "", err
	}

	userData, err := p.UserData()
	if err != nil {
		return "", err
	}

	unsigned = fmt.Sprintf("%s:%s", unsigned, userData)
	return unsigned, err
}
