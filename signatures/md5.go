package signatures

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(body []byte) string {
	bodyMD5 := md5.New()
	bodyMD5.Write(body)
	return hex.EncodeToString(bodyMD5.Sum(nil))
}
