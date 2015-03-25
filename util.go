package pusher

import (
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var ChannelValidationRegex = regexp.MustCompile("^[-a-zA-Z0-9_=@,.;]+$")

func authTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func parseAuthRequestParams(_params []byte) (string, string) {
	params, _ := url.ParseQuery(string(_params))
	return params["channel_name"][0], params["socket_id"][0]
}

func channelsAreValid(channels []string) bool {
	for _, channel := range channels {
		if len(channel) > 200 || !ChannelValidationRegex.MatchString(channel) {
			return false
		}
	}
	return true
}
