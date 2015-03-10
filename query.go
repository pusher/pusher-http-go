package pusher

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

const auth_version = "1.0"
const domain = "http://api.pusherapp.com"

type Query struct {
	request_method, path, key, secret string
	body                              []byte
}

func (q *Query) body_md5() string {
	_body_md5 := md5.New()
	_body_md5.Write([]byte(q.body))
	return hex.EncodeToString(_body_md5.Sum(nil))
}

func (q *Query) auth_timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func (q *Query) pre_signature_string() string {
	return "auth_key=" + q.key + "&" +
		"auth_timestamp=" + q.auth_timestamp() + "&" +
		"auth_version=" + auth_version + "&" +
		"body_md5=" + q.body_md5()
}

func (q *Query) sign() (string, string) {

	pre_signature_string := q.pre_signature_string()

	to_sign := q.request_method +
		"\n" +
		q.path +
		"\n" +
		pre_signature_string

	fmt.Println(to_sign)

	_auth_signature := hmac.New(sha256.New, []byte(q.secret))
	_auth_signature.Write([]byte(to_sign))
	return pre_signature_string, hex.EncodeToString(_auth_signature.Sum(nil))
}

func (q *Query) generate() string {
	pre_signature, auth_signature := q.sign()
	return domain + q.path + "?" + pre_signature + "&auth_signature=" + auth_signature
}
