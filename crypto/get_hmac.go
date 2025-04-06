package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GetHmac(message string, key string) string {
	hmac := hmac.New(sha256.New, []byte(key))
	hmac.Write([]byte(message))
	digest := hmac.Sum(nil)
	return base64.StdEncoding.EncodeToString(digest)
}
