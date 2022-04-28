package middleware

import (
	contracts_stores_basicauth "echo-starter/internal/contracts/stores/basicauth"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	core_wellknown "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/wellknown"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
)

type (
	// KeyAuthConfig defines the config for KeyAuth middleware.
	KeyAuthConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo_middleware.Skipper

		// KeyLookup is a string in the form of "<source>:<name>" or "<source>:<name>,<source>:<name>" that is used
		// to extract key from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>" or "header:<name>:<cut-prefix>"
		// 			`<cut-prefix>` is argument value to cut/trim prefix of the extracted value. This is useful if header
		//			value has static prefix like `Authorization: <auth-scheme> <authorisation-parameters>` where part that we
		//			want to cut is `<auth-scheme> ` note the space at the end.
		//			In case of basic authentication `Authorization: Basic <credentials>` prefix we want to remove is `Basic `.
		// - "query:<name>"
		// - "form:<name>"
		// - "cookie:<name>"
		// Multiple sources example:
		// - "header:Authorization,header:X-Api-Key"
		KeyLookup string

		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Bearer".
		AuthScheme string

		// Validator is a function to validate key.
		// Required.
		Validator KeyAuthValidator

		// ErrorHandler defines a function which is executed for an invalid key.
		// It may be used to define a custom error.
		ErrorHandler KeyAuthErrorHandler

		// ContinueOnIgnoredError allows the next middleware/handler to be called when ErrorHandler decides to
		// ignore the error (by returning `nil`).
		// This is useful when parts of your site/api allow public access and some authorized routes provide extra functionality.
		// In that case you can use ErrorHandler to set a default public key auth value in the request context
		// and continue. Some logic down the remaining execution chain needs to check that (public) key auth value then.
		ContinueOnIgnoredError bool
	}

	// KeyAuthValidator defines a function to validate KeyAuth credentials.
	KeyAuthValidator func(auth string, c echo.Context) (bool, error)

	// KeyAuthErrorHandler defines a function which is executed for an invalid key.
	KeyAuthErrorHandler func(err error, c echo.Context) error
)

var (
	// DefaultKeyAuthConfig is the default KeyAuth middleware config.
	DefaultKeyAuthConfig = KeyAuthConfig{
		Skipper:    echo_middleware.DefaultSkipper,
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
	}
)

// ErrKeyAuthMissing is error type when KeyAuth middleware is unable to extract value from lookups
type ErrKeyAuthMissing struct {
	Err error
}

// Error returns errors text
func (e *ErrKeyAuthMissing) Error() string {
	return e.Err.Error()
}

// Unwrap unwraps error
func (e *ErrKeyAuthMissing) Unwrap() error {
	return e.Err
}

// BasicAuthWithIBasicAuthStore returns an BasicAuth middleware.
//
// For valid credentials it calls the next handler.
// For missing or invalid credentials, it sends "401 - Unauthorized" response.
func KeyAuthWithIBasicAuthStore(container di.Container) echo.MiddlewareFunc {
	store := contracts_stores_basicauth.GetIBasicAuthStoreFromContainer(container)
	var validator = func(auth string, _ echo.Context) (bool, error) {
		return store.Validate("X-Api-Key", auth)
	}
	c := DefaultKeyAuthConfig
	c.KeyLookup = "header:X-Api-Key"
	c.Validator = validator
	return KeyAuthWithConfig(c)
}

// KeyAuth returns an KeyAuth middleware.
//
// For valid key it calls the next handler.
// For invalid key, it sends "401 - Unauthorized" response.
// For missing key, it sends "400 - Bad Request" response.
func KeyAuth(fn KeyAuthValidator) echo.MiddlewareFunc {
	c := DefaultKeyAuthConfig
	c.Validator = fn
	return KeyAuthWithConfig(c)
}

// KeyAuthWithConfig returns an KeyAuth middleware with config.
// See `KeyAuth()`.
func KeyAuthWithConfig(config KeyAuthConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultKeyAuthConfig.Skipper
	}
	// Defaults
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultKeyAuthConfig.AuthScheme
	}
	if config.KeyLookup == "" {
		config.KeyLookup = DefaultKeyAuthConfig.KeyLookup
	}
	if config.Validator == nil {
		panic("echo: key-auth middleware requires a validator function")
	}

	extractors, err := createExtractors(config.KeyLookup, config.AuthScheme)
	if err != nil {
		panic(err)
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
			for _, extractor := range extractors {
				keys, err := extractor(c)
				if err != nil {
					continue
				}
				for _, key := range keys {
					logger.Trace().Str("key", key).Msg("key-auth: key extracted")
					valid, err := config.Validator(key, c)
					if err != nil {
						logger.Trace().Str("key", key).Err(err).Msg("key-auth: key validation failed")
						continue
					}
					if valid {
						claimsPrincipal.AddClaim(contracts_core_claimsprincipal.Claim{
							Type:  core_wellknown.ClaimTypeAuthenticated,
							Value: "*"})
						claimsPrincipal.AddClaim(contracts_core_claimsprincipal.Claim{
							Type:  "auth_type",
							Value: "api-key"})
						return next(c)
					}
				}
			}

			return next(c)
		}
	}
}
