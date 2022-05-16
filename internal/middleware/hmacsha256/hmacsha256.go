package hmacsha256

import (
	"fmt"
	"net/http"

	"io/ioutil"

	"github.com/labstack/echo/v4"
)

func ValidateHMACSHA256() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			index := 0
			headers := c.Request().Header
			b, _ := ioutil.ReadAll(c.Request().Body)
			signature, _ := ComputeHash(b, []byte("secret"))
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
