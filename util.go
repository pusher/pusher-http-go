package pusher

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var channelValidationRegex = regexp.MustCompile("^[-a-zA-Z0-9_=@,.;]+$")
var socketIDValidationRegex = regexp.MustCompile(`\A\d+\.\d+\z`)
var maxChannelNameSize = 200

func authTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func parseAuthRequestParams(_params []byte) (channelName string, socketID string, err error) {
	params, err := url.ParseQuery(string(_params))
	if err != nil {
		return
	}
	if _, ok := params["channel_name"]; !ok {
		return "", "", errors.New("Channel param not found")
	}
	if _, ok := params["socket_id"]; !ok {
		return "", "", errors.New("Socket_id not found")
	}
	return params["channel_name"][0], params["socket_id"][0], nil
}

func validChannel(channel string) bool {
	if len(channel) > maxChannelNameSize || !channelValidationRegex.MatchString(channel) {
		return false
	}
	return true
}

func channelsAreValid(channels []string) bool {
	for _, channel := range channels {
		if !validChannel(channel) {
			return false
		}
	}
	return true
}

func isEncryptedChannel(channel string) bool {
	if strings.HasPrefix(channel, "private-encrypted-") {
		return true
	}
	return false
}

func validEncryptionKey(encryptionKey string) bool {
	return len(encryptionKey) == 32
}

func validateSocketID(socketID *string) (err error) {
	if (socketID == nil) || socketIDValidationRegex.MatchString(*socketID) {
		return
	}
	return errors.New("socket_id invalid")
}
