package hmacsha256

import (
	"fmt"
	"net/http"

	"io/ioutil"

	"github.com/labstack/echo/v4"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	core_wellknown "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/wellknown"

	di "github.com/fluffy-bunny/sarulabsdi"
)

const middlewareLogName = "hmac-sha256"

func ValidateHMACSHA256() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			scopedContainer := c.Get(core_wellknown.SCOPED_CONTAINER_KEY).(di.Container)
			loggerObj := contracts_logger.GetILoggerFromContainer(scopedContainer)
			logger := loggerObj.GetLogger().With().Caller().Str("middleware", middlewareLogName).Logger()

			index := 0
			headers := c.Request().Header
			b, _ := ioutil.ReadAll(c.Request().Body)
			signature, _ := ComputeHash(b, []byte("secret"))
			logger = logger.With().Str("signature", signature).Logger()

			found := false
			for {
				key := fmt.Sprintf("X-Mapped-Signature-%d", index)
				value, ok := headers[key]
				if ok {
					logger.Debug().
						Str("signature-key", key).
						Str("signature-value", value[0]).Msg("found")
					found = true
					if value[0] == signature {
						logger.Debug().Msg("signature-valid")
						return next(c)
					}
				}
				break
			}
			if found {
				logger.Error().Msg("hmac sha256 validation failed")
				return c.JSON(http.StatusUnauthorized, "hmac sha256 validation failed")
			}

			return next(c)
		}
	}
}
