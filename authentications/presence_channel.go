package authentications

import (
	"fmt"
)

type PresenceChannel struct {
	Body []byte
	Member
}

func (p *PresenceChannel) StringToSign() (unsigned string, err error) {
	privateChannelRequest := &PrivateChannel{p.Body}
	if unsigned, err = privateChannelRequest.StringToSign(); err != nil {
		return
	}
	var userData string
	if userData, err = p.UserData(); err != nil {
		return
	}

	unsigned = fmt.Sprintf("%s:%s", unsigned, userData)
	return
}
