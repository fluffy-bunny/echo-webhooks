package hmacsha256

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func ComputeHash(msg, key []byte) (base64signature string, err error) {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	encodedString := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return encodedString, nil
}
