package pusher

import (
	"net/url"
	"regexp"
	"strconv"
	"time"
)

//this is Ruby :(
func auth_timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func parseAuthRequestParams(_params []byte) (string, string) {
	params, _ := url.ParseQuery(string(_params))
	return params["channel_name"][0], params["socket_id"][0]
}

func channelsAreValid(channels []string) bool {
	r, _ := regexp.Compile("^[-a-zA-Z0-9_=@,.;]+$")
	for _, channel := range channels {
		if !r.MatchString(channel) || len(channel) > 200 {
			return false
		}
	}
	return true
}
