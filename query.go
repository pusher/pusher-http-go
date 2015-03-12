package pusher

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const auth_version = "1.0"
const domain = "http://api.pusherapp.com"

type Query struct {
	request_method, path, key, secret string
	body                              []byte
	additional_queries                map[string]string
}

func (q *Query) body_md5() string {
	_body_md5 := md5.New()
	_body_md5.Write([]byte(q.body))
	return hex.EncodeToString(_body_md5.Sum(nil))
}

func (q *Query) auth_timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func (q *Query) unsigned_params() url.Values {

	params := url.Values{}
	params.Add("auth_key", q.key)
	params.Add("auth_timestamp", q.auth_timestamp())
	params.Add("auth_version", auth_version)
	params.Add("body_md5", q.body_md5())

	if q.additional_queries != nil {
		for key, value := range q.additional_queries {
			params.Add(key, value)
		}
	}

	return params
}

func (q *Query) sign() (url.Values, string) {

	unsigned_params := q.unsigned_params()

	to_sign := q.request_method +
		"\n" +
		q.path +
		"\n" +
		unsigned_params.Encode()

	_auth_signature := hmac.New(sha256.New, []byte(q.secret))
	_auth_signature.Write([]byte(to_sign))
	return unsigned_params, hex.EncodeToString(_auth_signature.Sum(nil))
}

func (q *Query) generate() string {
	params, auth_signature := q.sign()
	params.Add("auth_signature", auth_signature)

	endpoint, _ := url.Parse(domain + q.path)
	endpoint.RawQuery = params.Encode()

	fmt.Print(endpoint.String())

	return endpoint.String()
}
