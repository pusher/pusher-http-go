package authentications

import (
	"fmt"
	"github.com/pusher/pusher-http-go/errors"
	"github.com/pusher/pusher-http-go/validate"
	"net/url"
)

type PrivateChannel struct {
	Body []byte
}

func (p *PrivateChannel) StringToSign() (string, error) {
	params, err := url.ParseQuery(string(p.Body))
	if err != nil {
		return "", err
	}

	channelNameWrapper, keyExists := params["channel_name"]
	if !keyExists || len(channelNameWrapper) == 0 {
		return "", errors.New("Channel param not found")
	}

	socketIDWrapper, keyExists := params["socket_id"]
	if !keyExists || len(socketIDWrapper) == 0 {
		return "", errors.New("Socket_id not found")
	}

	channelName := channelNameWrapper[0]
	if len(channelName) == 0 {
		return "", errors.New("Channel name cannot be blank")
	}

	socketID := socketIDWrapper[0]
	if err = validate.SocketID(&socketID); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", socketID, channelName), nil
}

func (p *PrivateChannel) UserData() (userData string, err error) {
	return
}
