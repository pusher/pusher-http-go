package signatures

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HMAC(unsigned, secret string) string {
	return hex.EncodeToString(signBytesHMAC([]byte(unsigned), []byte(secret)))
}

func CheckHMAC(result, secret string, body []byte) bool {
	expected := signBytesHMAC(body, []byte(secret))
	resultBytes, err := hex.DecodeString(result)
	if err != nil {
		return false
	}

	return hmac.Equal(expected, resultBytes)
}

func signBytesHMAC(unsigned, secret []byte) []byte {
	signature := hmac.New(sha256.New, secret)
	signature.Write(unsigned)
	return signature.Sum(nil)
}
