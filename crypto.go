package pusher

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	// "fmt"
	"strings"
)

func hmacSignature(toSign, secret string) string {
	return hex.EncodeToString(hmacBytes([]byte(toSign), []byte(secret)))
}

func hmacBytes(toSign, secret []byte) []byte {
	_authSignature := hmac.New(sha256.New, secret)
	_authSignature.Write(toSign)
	return _authSignature.Sum(nil)
}

func checkSignature(result, secret string, body []byte) bool {
	expected := hmacBytes(body, []byte(secret))
	resultBytes, err := hex.DecodeString(result)

	if err != nil {
		return false
	}
	return hmac.Equal(expected, resultBytes)
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
