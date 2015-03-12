package pusher

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
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

	params := url.Values{
		"auth_key":       {q.key},
		"auth_timestamp": {q.auth_timestamp()},
		"auth_version":   {auth_version},
	}

	if q.body != nil {
		params.Add("body_md5", q.body_md5())
	}

	if q.additional_queries != nil {
		for key, values := range q.additional_queries {
			params.Add(key, values)
		}
	}

	return params
}

func (q *Query) sign(params url.Values) string {

	to_sign := q.request_method +
		"\n" +
		q.path +
		"\n" +
		q.unescape_query(params)

	_auth_signature := hmac.New(sha256.New, []byte(q.secret))
	_auth_signature.Write([]byte(to_sign))
	return hex.EncodeToString(_auth_signature.Sum(nil))
}

func (q *Query) unescape_query(query url.Values) string {
	unesc, _ := url.QueryUnescape(query.Encode())
	return unesc
}

func (q *Query) generate() string {

	params := q.unsigned_params()

	auth_signature := q.sign(params)
	params.Add("auth_signature", auth_signature)

	endpoint, _ := url.Parse(domain + q.path)
	endpoint.RawQuery = q.unescape_query(params)

	return endpoint.String()
}
