package requests

import (
	"github.com/pusher/pusher/signatures"
	"net/url"
)

type Params struct {
	Body    []byte
	Queries map[string]string
	Channel string
}

const authVersion = "1.0"

func (p *Params) URLValues(key string) (values *url.Values) {
	values = &url.Values{
		"auth_key":       {key},
		"auth_timestamp": {authClock.Now()},
		"auth_version":   {authVersion},
	}

	if p.Body != nil {
		values.Add("body_md5", signatures.MD5(p.Body))
	}

	if p.Queries != nil {
		for key, value := range p.Queries {
			values.Add(key, value)
		}
	}
	return
}
