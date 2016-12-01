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

func requestURL(uc *urlConfig, request *requests.Request, params *requests.Params) (*url.URL, error) {
	values := params.URLValues(uc.key)

	var path string
	if params.Channel != "" {
		path = fmt.Sprintf(request.PathPattern, uc.appID, params.Channel)
	} else {
		path = fmt.Sprintf(request.PathPattern, uc.appID)
	}

	encodedURLValues := values.Encode()
	urlUnescaped, err := url.QueryUnescape(encodedURLValues)
	if err != nil {
		err = errors.New(fmt.Sprintf("%s could not be unescaped - %v", encodedURLValues, err))
		return nil, err
	}

	unsigned := s.Join([]string{request.Method, path, urlUnescaped}, "\n")
	signed := signatures.HMAC(unsigned, uc.secret)
	values.Add("auth_signature", signed)

	u := &url.URL{
		Scheme:   uc.scheme,
		Host:     uc.host,
		Path:     path,
		RawQuery: values.Encode(),
	}

	return u, nil
}
