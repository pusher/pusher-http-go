package pusher

import (
	"net/url"
	"strconv"
	"time"
)

func auth_timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func parseAuthRequestParams(_params []byte) (string, string) {
	params, _ := url.ParseQuery(string(_params))
	return params["channel_name"][0], params["socket_id"][0]
}
