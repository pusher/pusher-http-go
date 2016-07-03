package authentications

import (
	"fmt"
	"github.com/pusher/pusher/errors"
	"github.com/pusher/pusher/validate"
	"net/url"
)

type PrivateChannel struct {
	Body []byte
}

func (p *PrivateChannel) StringToSign() (unsigned string, err error) {
	var (
		params             url.Values
		keyExists          bool
		channelNameWrapper []string
		socketIDWrapper    []string
		channelName        string
		socketID           string
	)

	if params, err = url.ParseQuery(string(p.Body)); err != nil {
		return
	}

	if channelNameWrapper, keyExists = params["channel_name"]; !keyExists || len(channelNameWrapper) == 0 {
		err = errors.New("Channel param not found")
		return
	}

	if socketIDWrapper, keyExists = params["socket_id"]; !keyExists || len(socketIDWrapper) == 0 {
		err = errors.New("Socket_id not found")
		return
	}

	if channelName = channelNameWrapper[0]; len(channelName) == 0 {
		err = errors.New("Channel name cannot be blank")
		return
	}

	socketID = socketIDWrapper[0]
	if err = validate.SocketID(&socketID); err != nil {
		return
	}

	unsigned = fmt.Sprintf("%s:%s", socketID, channelName)
	return
}

func (p *PrivateChannel) UserData() (userData string, err error) {
	return
}
