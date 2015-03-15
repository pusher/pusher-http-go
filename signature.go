package pusher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HMACSignature(to_sign, secret string) string {
	_auth_signature := hmac.New(sha256.New, []byte(secret))
	_auth_signature.Write([]byte(to_sign))
	return hex.EncodeToString(_auth_signature.Sum(nil))
}

func CheckSignature(result, body, secret string) bool {
	expected := HMACSignature(body, secret)
	return result == expected
}
