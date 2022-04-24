//go:build go1.15
// +build go1.15

package middleware

import (
	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	core_contracts_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oidc"
	core_wellknown "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/wellknown"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	jose_jwt "gopkg.in/square/go-jose.v2/jwt"
)

type (
	// JWTConfig defines the config for JWT middleware.
	JWTConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo_middleware.Skipper

		// BeforeFunc defines a function which is executed just before the middleware.
		BeforeFunc echo_middleware.BeforeFunc

		// SuccessHandler defines a function which is executed for a valid token before middleware chain continues with next
		// middleware or handler.
		SuccessHandler JWTSuccessHandler

		// ErrorHandler defines a function which is executed for an invalid token.
		// It may be used to define a custom JWT error.
		ErrorHandler JWTErrorHandler

		// ErrorHandlerWithContext is almost identical to ErrorHandler, but it's passed the current context.
		ErrorHandlerWithContext JWTErrorHandlerWithContext

		// ContinueOnIgnoredError allows the next middleware/handler to be called when ErrorHandlerWithContext decides to
		// ignore the error (by returning `nil`).
		// This is useful when parts of your site/api allow public access and some authorized routes provide extra functionality.
		// In that case you can use ErrorHandlerWithContext to set a default public JWT token value in the request context
		// and continue. Some logic down the remaining execution chain needs to check that (public) token value then.
		ContinueOnIgnoredError bool

		// TokenLookup is a string in the form of "<source>:<name>" or "<source>:<name>,<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>" or "header:<name>:<cut-prefix>"
		// 			`<cut-prefix>` is argument value to cut/trim prefix of the extracted value. This is useful if header
		//			value has static prefix like `Authorization: <auth-scheme> <authorisation-parameters>` where part that we
		//			want to cut is `<auth-scheme> ` note the space at the end.
		//			In case of JWT tokens `Authorization: Bearer <token>` prefix we cut is `Bearer `.
		// If prefix is left empty the whole value is returned.
		// - "query:<name>"
		// - "param:<name>"
		// - "cookie:<name>"
		// - "form:<name>"
		// Multiple sources example:
		// - "header:Authorization,cookie:myowncookie"
		TokenLookup string

		// TokenLookupFuncs defines a list of user-defined functions that extract JWT token from the given context.
		// This is one of the two options to provide a token extractor.
		// The order of precedence is user-defined TokenLookupFuncs, and TokenLookup.
		// You can also provide both if you want.
		TokenLookupFuncs []ValuesExtractor

		// AuthScheme to be used in the Authorization header.
		// Optional. Default value "Bearer".
		AuthScheme string
	}

	// JWTSuccessHandler defines a function which is executed for a valid token.
	JWTSuccessHandler func(c echo.Context)

	// JWTErrorHandler defines a function which is executed for an invalid token.
	JWTErrorHandler func(err error) error

	// JWTErrorHandlerWithContext is almost identical to JWTErrorHandler, but it's passed the current context.
	JWTErrorHandlerWithContext func(err error, c echo.Context) error
)

var (
	// DefaultJWTConfig is the default JWT auth middleware config.
	DefaultJWTConfig = JWTConfig{
		Skipper:          echo_middleware.DefaultSkipper,
		TokenLookup:      "header:" + echo.HeaderAuthorization,
		TokenLookupFuncs: nil,
		AuthScheme:       "Bearer",
	}
)

// JWT returns a JSON Web Token (JWT) auth middleware.
//
// For valid token, it sets the user in context and calls next handler.
// For invalid token, it returns "401 - Unauthorized" error.
// For missing token, it returns "400 - Bad Request" error.
//
// See: https://jwt.io/introduction
// See `JWTConfig.TokenLookup`
func JWT(root di.Container) echo.MiddlewareFunc {
	c := DefaultJWTConfig
	return JWTWithConfig(root, c)
}

const middlewareLogName = "jwt-to-claims-principal"

// JWTWithConfig returns a JWT auth middleware with config.
// See: `JWT()`.
func JWTWithConfig(root di.Container, config JWTConfig) echo.MiddlewareFunc {
	log.Info().Msg("JWT to Claims Middleware")
	oidcAuthenticator, _ := core_contracts_oidc.SafeGetIOIDCAuthenticatorFromContainer(root)
	if oidcAuthenticator == nil {
		log.Info().Msg("JWT to Claims Middleware: OIDC Authenticator not found in container")
	}
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultJWTConfig.Skipper
	}

	if config.TokenLookup == "" && len(config.TokenLookupFuncs) == 0 {
		config.TokenLookup = DefaultJWTConfig.TokenLookup
	}
	if config.AuthScheme == "" {
		config.AuthScheme = DefaultJWTConfig.AuthScheme
	}

	extractors, err := CreateExtractors(config.TokenLookup, config.AuthScheme)
	if err != nil {
		panic(err)
	}
	if len(config.TokenLookupFuncs) > 0 {
		extractors = append(config.TokenLookupFuncs, extractors...)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			log.Trace().Msg("JWT to Claims Middleware - ENTER")
			defer log.Trace().Msg("JWT to Claims Middleware - EXIT")
			if config.Skipper(c) {
				return next(c)
			}

			if config.BeforeFunc != nil {
				config.BeforeFunc(c)
			}

			scopedContainer := c.Get(core_wellknown.SCOPED_CONTAINER_KEY).(di.Container)

			loggerObj := contracts_logger.GetILoggerFromContainer(scopedContainer)
			logger := loggerObj.GetLogger().With().Str("middleware", middlewareLogName).Logger()

			claimsPrincipal := contracts_core_claimsprincipal.GetIClaimsPrincipalFromContainer(scopedContainer)

			for _, extractor := range extractors {
				auths, err := extractor(c)
				if err != nil {
					continue
				}
				for _, auth := range auths {
					logger.Trace().Str("token", auth).Send()
					if ok, _ := isJWT(auth); ok {
						if oidcAuthenticator != nil {

							accessToken, err := oidcAuthenticator.ValidateJWTAccessToken(auth)
							if err != nil {
								logger.Error().Err(err).Msg("ValidateJWTAccessToken failed")
								continue
							}
							accessTokenClaims := accessToken.ToClaims()
							for _, claim := range accessTokenClaims {
								claimsPrincipal.AddClaim(*claim)
							}
							claimsPrincipal.AddClaim(contracts_core_claimsprincipal.Claim{
								Type:  core_wellknown.ClaimTypeAuthenticated,
								Value: "*"})
							logger.Trace().Interface("claims", claimsPrincipal.GetClaims()).Send()
						}
					}
				}
			}

			return next(c)

		}
	}
}

func isJWT(token string) (bool, error) {
	_, err := jose_jwt.ParseSigned(token)
	if err != nil {
		return false, err
	}
	return true, nil
}
