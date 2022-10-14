package pusher

import (
	"errors"
	"fmt"
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

func parseUserAuthenticationRequestParams(_params []byte) (socketID string, err error) {
	params, err := url.ParseQuery(string(_params))
	if err != nil {
		return
	}
	if _, ok := params["socket_id"]; !ok {
		return "", errors.New("socket_id not found")
	}
	return params["socket_id"][0], nil
}

func parseChannelAuthorizationRequestParams(_params []byte) (channelName string, socketID string, err error) {
	params, err := url.ParseQuery(string(_params))
	if err != nil {
		return
	}
	if _, ok := params["channel_name"]; !ok {
		return "", "", errors.New("channel_name not found")
	}
	if _, ok := params["socket_id"]; !ok {
		return "", "", errors.New("socket_id not found")
	}
	return params["channel_name"][0], params["socket_id"][0], nil
}

func validUserId(userId string) bool {
	length := len(userId)
	return length > 0 && length < maxChannelNameSize
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

func validateUserData(userData map[string]interface{}) (err error) {
	_id, ok := userData["id"]
	if !ok || _id == nil {
		return errors.New("Missing id in user data")
	}
	var id string
	id, ok = _id.(string)
	if !ok {
		return errors.New("id field in user data is not a string")
	}
	if !validUserId(id) {
		return fmt.Errorf("Invalid id in user data: '%s'", id)
	}
	return
}

func validateSocketID(socketID *string) (err error) {
	if (socketID == nil) || socketIDValidationRegex.MatchString(*socketID) {
		return
	}
	return errors.New("socket_id invalid")
}
