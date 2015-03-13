package pusher

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"

	"strconv"
	"time"
)

const auth_version = "1.0"
const domain = "http://api.pusherapp.com"

type Url struct {
	request_method, path, key, secret string
	body                              []byte
	additional_queries                map[string]string
}

func (u *Url) body_md5() string {
	_body_md5 := md5.New()
	_body_md5.Write([]byte(u.body))
	return hex.EncodeToString(_body_md5.Sum(nil))
}

func (u *Url) auth_timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func (u *Url) unsigned_params() url.Values {

	params := url.Values{
		"auth_key":       {u.key},
		"auth_timestamp": {u.auth_timestamp()},
		"auth_version":   {auth_version},
	}

	if u.body != nil {
		params.Add("body_md5", u.body_md5())
	}

	if u.additional_queries != nil {
		for key, values := range u.additional_queries {
			params.Add(key, values)
		}
	}

	return params
}

func (u *Url) unescape_url(Url url.Values) string {
	unesc, _ := url.QueryUnescape(Url.Encode())
	return unesc
}

func (u *Url) generate() string {

	params := u.unsigned_params()

	string_to_sign := u.request_method + "\n" + u.path + "\n" + u.unescape_url(params)

	auth_signature := HMACSignature(string_to_sign, u.secret)
	params.Add("auth_signature", auth_signature)

	endpoint, _ := url.Parse(domain + u.path)
	endpoint.RawQuery = u.unescape_url(params)

	return endpoint.String()
}
