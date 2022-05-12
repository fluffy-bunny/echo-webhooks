package middleware

import (
	contracts_stores_basicauth "echo-starter/internal/contracts/stores/basicauth"
	"encoding/base64"
	"strings"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	core_wellknown "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/wellknown"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

type (
	// BasicAuthConfig defines the config for BasicAuth middleware.
	BasicAuthConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo_middleware.Skipper

		// Validator is a function to validate BasicAuth credentials.
		// Required.
		Validator BasicAuthValidator

		// Realm is a string to define realm attribute of BasicAuth.
		// Default value "Restricted".
		Realm string
	}

	// BasicAuthValidator defines a function to validate BasicAuth credentials.
	BasicAuthValidator func(string, string, echo.Context) (bool, error)
)

const (
	basic        = "basic"
	defaultRealm = "Restricted"
)

var (
	// DefaultBasicAuthConfig is the default BasicAuth middleware config.
	DefaultBasicAuthConfig = BasicAuthConfig{
		Skipper: echo_middleware.DefaultSkipper,
		Realm:   defaultRealm,
	}
)

// BasicAuthWithIBasicAuthStore returns an BasicAuth middleware.
//
// For valid credentials it calls the next handler.
// For missing or invalid credentials, it sends "401 - Unauthorized" response.
func BasicAuthWithIBasicAuthStore(container di.Container) echo.MiddlewareFunc {
	store := contracts_stores_basicauth.GetIBasicAuthStoreFromContainer(container)
	var validator = func(username string, password string, _ echo.Context) (bool, error) {
		return store.Validate(username, password)
	}
	c := DefaultBasicAuthConfig
	c.Validator = validator
	return BasicAuthWithConfig(c)
}

// BasicAuth returns an BasicAuth middleware.
//
// For valid credentials it calls the next handler.
// For missing or invalid credentials, it sends "401 - Unauthorized" response.
func BasicAuth(fn BasicAuthValidator) echo.MiddlewareFunc {
	c := DefaultBasicAuthConfig
	c.Validator = fn
	return BasicAuthWithConfig(c)
}

// BasicAuthWithConfig returns an BasicAuth middleware with config.
// See `BasicAuth()`.
func BasicAuthWithConfig(config BasicAuthConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Validator == nil {
		panic("echo: basic-auth middleware requires a validator function")
	}
	if config.Skipper == nil {
		config.Skipper = DefaultBasicAuthConfig.Skipper
	}
	if config.Realm == "" {
		config.Realm = defaultRealm
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}
			scopedContainer := c.Get(core_wellknown.SCOPED_CONTAINER_KEY).(di.Container)

			loggerObj := contracts_logger.GetILoggerFromContainer(scopedContainer)
			logger := loggerObj.GetLogger().With().Str("middleware", middlewareLogName).Logger()

			claimsPrincipal := contracts_core_claimsprincipal.GetIClaimsPrincipalFromContainer(scopedContainer)
			auth := c.Request().Header.Get(echo.HeaderAuthorization)
			l := len(basic)

			if len(auth) > l+1 && strings.EqualFold(auth[:l], basic) {
				for {
					logger.Trace().Msg("basic auth header found")
					b, err := base64.StdEncoding.DecodeString(auth[l+1:])
					if err != nil {
						break
					}
					cred := string(b)
					for i := 0; i < len(cred); i++ {
						if cred[i] == ':' {
							// Verify credentials
							valid, err := config.Validator(cred[:i], cred[i+1:], c)
							if err != nil {
								break
							} else if valid {
								claimsPrincipal.AddClaim(contracts_core_claimsprincipal.Claim{
									Type:  core_wellknown.ClaimTypeAuthenticated,
									Value: "*"})
								claimsPrincipal.AddClaim(contracts_core_claimsprincipal.Claim{
									Type:  "auth_type",
									Value: "basic"})
							}
							break
						}
					}
					break
				}
			}

			return next(c)
		}
	}
}
