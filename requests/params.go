package requests

import (
	"github.com/pusher/pusher-http-go/signatures"
	"net/url"
)

const authVersion = "1.0"

type Params struct {
	Body    []byte
	Queries map[string]string
	Channel string
}

func (p *Params) URLValues(key string) *url.Values {
	values := &url.Values{
		"auth_key":       {key},
		"auth_timestamp": {authClock.Now()},
		"auth_version":   {authVersion},
	}

	if p.Body != nil {
		values.Add("body_md5", signatures.MD5(p.Body))
	}

	if p.Queries != nil {
		for k, v := range p.Queries {
			values.Add(k, v)
		}
	}

	return values
}
