package pusher

import (
	"fmt"
	"github.com/pusher/pusher-http-go/errors"
	"github.com/pusher/pusher-http-go/requests"
	"github.com/pusher/pusher-http-go/signatures"
	"net/url"
	s "strings"
)

type urlConfig struct {
	appID, key, secret, host, scheme string
}

func requestURL(p *urlConfig, request *requests.Request, params *requests.Params) (u *url.URL, err error) {
	values := params.URLValues(p.key)

	var path string
	if params.Channel != "" {
		path = fmt.Sprintf(request.PathPattern, p.appID, params.Channel)
	} else {
		path = fmt.Sprintf(request.PathPattern, p.appID)
	}

	var urlUnescaped string
	encodedURLValues := values.Encode()
	if urlUnescaped, err = url.QueryUnescape(encodedURLValues); err != nil {
		err = errors.New(fmt.Sprintf("%s could not be unescaped - %v", encodedURLValues, err))
		return
	}

	unsigned := s.Join([]string{request.Method, path, urlUnescaped}, "\n")
	signed := signatures.HMAC(unsigned, p.secret)
	values.Add("auth_signature", signed)

	u = &url.URL{
		Scheme:   p.scheme,
		Host:     p.host,
		Path:     path,
		RawQuery: values.Encode(),
	}

	return
}
