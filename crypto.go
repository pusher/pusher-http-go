package pusher

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func hmacSignature(toSign, secret string) string {
	_authSignature := hmac.New(sha256.New, []byte(secret))
	_authSignature.Write([]byte(toSign))
	return hex.EncodeToString(_authSignature.Sum(nil))
}

func checkSignature(result, body, secret string) bool {
	expected := hmacSignature(body, secret)
	return result == expected
}

func createAuthMap(key, secret, stringToSign string) map[string]string {
	authSignature := hmacSignature(stringToSign, secret)
	authString := strings.Join([]string{key, authSignature}, ":")
	return map[string]string{"auth": authString}
}

func md5Signature(body []byte) string {
	_bodyMD5 := md5.New()
	_bodyMD5.Write([]byte(body))
	return hex.EncodeToString(_bodyMD5.Sum(nil))
}
