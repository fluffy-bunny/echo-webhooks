package hmacsha256

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"io/ioutil"

	"github.com/labstack/echo/v4"
)

const middlewareLogName = "hmacsha256"

func signMessageHMAC256(msg, key []byte) (base64signature string, err error) {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	encodedString := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return encodedString, nil
}
func ValidateHMACSHA256() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			index := 0
			headers := c.Request().Header
			b, _ := ioutil.ReadAll(c.Request().Body)
			signature, _ := signMessageHMAC256(b, []byte("secret333"))
			found := false
			for {
				key := fmt.Sprintf("X-Mapped-Signature-%d", index)
				value, ok := headers[key]
				if ok {
					found = true
					if value[0] == signature {
						return next(c)
					}
				}
				break
			}
			if found {
				return c.JSON(http.StatusUnauthorized, "hmac sha256 validation failed")
			}

			return next(c)
		}
	}
}
