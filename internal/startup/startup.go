package startup

import (
	"context"
	echostarter_auth "echo-starter/internal/auth"
	contracts_config "echo-starter/internal/contracts/config"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	core_echo "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo"
	core_utils "github.com/fluffy-bunny/grpcdotnetgo/pkg/utils"
	"github.com/quasoft/memstore"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"

	services_handlers_about "echo-starter/internal/services/handlers/about"
	app_session "echo-starter/internal/session"

	"net/http"

	"github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/startup"
	echo_contracts_startup "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/startup"

	"github.com/gorilla/securecookie"

	services_auth_cookie_token_store "echo-starter/internal/services/auth/cookie_token_store"

	services_auth_session_token_store "echo-starter/internal/services/auth/session_token_store"
	services_handlers_api_webhook "echo-starter/internal/services/handlers/api/webhook"
	services_handlers_api_webhookapikey "echo-starter/internal/services/handlers/api/webhookapikey"
	services_handlers_api_webhookbasicauth "echo-starter/internal/services/handlers/api/webhookbasicauth"
	services_handlers_api_webhooknoauth "echo-starter/internal/services/handlers/api/webhooknoauth"

	services_stores_basicauth "echo-starter/internal/services/stores/basicauth"

	services_handlers_channel "echo-starter/internal/services/handlers/channel"
	services_handlers_healthz "echo-starter/internal/services/handlers/healthz"
	services_handlers_ready "echo-starter/internal/services/handlers/ready"
	services_probes_database "echo-starter/internal/services/probes/database"
	services_probes_oidc "echo-starter/internal/services/probes/oidc"

	// ACCOUNTS
	//----------------------------------------------------------------------------------------------------------------------

	core_contracts_oauth2 "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oauth2"

	core_contracts_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/oidc"
	core_services_oidc "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/oidc"

	core_services_oauth2_github "github.com/fluffy-bunny/grpcdotnetgo/pkg/services/oauth2/github"

	core_contracts_session "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/session"
	core_middleware_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/middleware/claimsprincipal"

	echo_middleware "echo-starter/internal/echo/middleware"
	//middleware_claimsprincipal "echo-starter/internal/middleware/claimsprincipal"

	contracts_sse "echo-starter/internal/contracts/sse"
	services_claimsprovider "echo-starter/internal/services/claimsprovider"
	services_sse "echo-starter/internal/services/sse"

	services_handlers_error "echo-starter/internal/services/handlers/error"
	services_handlers_home "echo-starter/internal/services/handlers/home"

	core_contracts "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/core"
	contracts_cookies "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/cookies"
	core_middleware_session "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/middleware/session"
	di "github.com/fluffy-bunny/sarulabsdi"
	redis "github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	redisstore "github.com/rbcervilla/redisstore/v8"
)

type Startup struct {
	echo_contracts_startup.CommonStartup
	config    *contracts_config.Config
	ctrl      *gomock.Controller
	container di.Container
	sseServer contracts_sse.IServerSideEventServer
}

func assertImplementation() {
	var _ echo_contracts_startup.IStartup = (*Startup)(nil)
}

func NewStartup() echo_contracts_startup.IStartup {
	startup := &Startup{
		config: &contracts_config.Config{},
		ctrl:   gomock.NewController(nil),
	}
	hooks := &echo_contracts_startup.Hooks{
		PostBuildHook: func(container di.Container) error {
			startup.container = container
			if startup.config.ApplicationEnvironment == "Development" {
				di.Dump(container)
			}
			return nil
		},
		PreStartHook: func(echo *echo.Echo) error {
			log.Info().Msgf("Starting %s", log.Logger.GetLevel().String())
			startup.sseServer = contracts_sse.GetIServerSideEventServerFromContainer(startup.container)
			return nil
		},
		PreShutdownHook: func(echo *echo.Echo) error {
			startup.sseServer.Shutdown()
			return nil
		},
	}

	startup.AddHooks(hooks)
	return startup
}

func (s *Startup) getSessionStore() sessions.Store {

	hashKey, err := base64.StdEncoding.DecodeString(s.config.SecureCookieHashKey)
	if err != nil {
		panic(err)
	}
	encryptionKey, err := base64.StdEncoding.DecodeString(s.config.SecureCookieEncryptionKey)
	if err != nil {
		panic(err)
	}

	switch s.config.SessionEngine {
	case "cookie":
		store := sessions.NewCookieStore(hashKey, encryptionKey)
		store.Options.Secure = true
		store.Options.HttpOnly = true
		store.Options.SameSite = http.SameSiteStrictMode
		store.Options.MaxAge = s.config.SessionMaxAgeSeconds
		return store
	case "inmemory":
		store := memstore.NewMemStore(hashKey, encryptionKey)
		store.Options.Secure = true
		store.Options.HttpOnly = true
		store.Options.SameSite = http.SameSiteStrictMode
		store.Options.MaxAge = s.config.SessionMaxAgeSeconds
		return store
	case "redis":
		client := redis.NewClient(&redis.Options{
			Addr:     s.config.RedisUrl,
			Password: s.config.RedisPassword,
		})

		// New default RedisStore
		store, err := redisstore.NewRedisStore(context.Background(), client)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create redis store")

		}
		store.Options(sessions.Options{
			Path:     "/",
			MaxAge:   s.config.SessionMaxAgeSeconds,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
		return store
	default:
		return nil
	}
}
func (s *Startup) RegisterStaticRoutes(e *echo.Echo) error {
	e.Static("/static", "./static")
	return nil
}

func (s *Startup) GetOptions() *startup.Options {
	return &startup.Options{
		Port: s.config.Port,
	}
}

func (s *Startup) GetConfigOptions() *core_contracts.ConfigOptions {
	prettyLog, err := strconv.ParseBool(os.Getenv("PRETTY_LOG"))
	if err != nil {
		prettyLog = false
	}

	return &core_contracts.ConfigOptions{
		RootConfig:             []byte(contracts_config.ConfigDefaultJSON),
		Destination:            s.config,
		LogLevel:               os.Getenv("LOG_LEVEL"),
		PrettyLog:              prettyLog,
		ApplicationEnvironment: os.Getenv("APPLICATION_ENVIRONMENT"),
	}
}
func (s *Startup) addSecureCookieOptions(builder *di.Builder) {
	// map our config to accessor funcs that other services need
	// SECURE COOKIE
	if core_utils.IsEmptyOrNil(s.config.SecureCookieHashKey) {
		fmt.Println("WARNING: SECURE_COOKIE_HASH_KEY must be set for production......")
		key := securecookie.GenerateRandomKey(32)
		encodedString := base64.StdEncoding.EncodeToString(key)
		s.config.SecureCookieHashKey = encodedString
		fmt.Printf("SECURE_COOKIE_HASH_KEY: %v\n", s.config.SecureCookieHashKey)
	}
	if core_utils.IsEmptyOrNil(s.config.SecureCookieEncryptionKey) {
		fmt.Println("WARNING: SECURE_COOKIE_ENCRYPTION_KEY must be set for production......")
		key := securecookie.GenerateRandomKey(32)
		encodedString := base64.StdEncoding.EncodeToString(key)
		s.config.SecureCookieEncryptionKey = encodedString
		fmt.Printf("SECURE_COOKIE_ENCRYPTION_KEY: %v\n", s.config.SecureCookieEncryptionKey)
	}

	contracts_cookies.AddSecureCookieConfigAccessorFunc(builder, func() *contracts_cookies.SecureCookieConfig {
		return &contracts_cookies.SecureCookieConfig{
			SecureCookieHashKey:       s.config.SecureCookieHashKey,
			SecureCookieEncryptionKey: s.config.SecureCookieEncryptionKey,
		}
	})
}
func (s *Startup) addAuthServices(builder *di.Builder) {

	switch s.config.AuthProvider {
	case "oidc":
		// AUTH SERVICES
		//----------------------------------------------------------------------------------------------------------------------
		core_contracts_oidc.AddGetOIDCAuthenticatorConfigFunc(builder, func() *core_contracts_oidc.AuthenticatorConfig {

			return &core_contracts_oidc.AuthenticatorConfig{
				Domain:       s.config.OIDC.Domain,
				ClientID:     s.config.OIDC.ClientID,
				ClientSecret: s.config.OIDC.ClientSecret,
				CallbackURL:  s.config.OIDC.CallbackURL,
				Insecure:     s.config.OIDC.Insecure,
			}
		})
		core_services_oidc.AddSingletonIOIDCAuthenticator(builder)

	case "github":
		core_contracts_oauth2.AddGetOAuth2AuthenticatorConfigFunc(builder, func() *oauth2.Config {
			return &oauth2.Config{
				ClientID:     s.config.OAuth2.ClientID,
				ClientSecret: s.config.OAuth2.ClientSecret,
				RedirectURL:  s.config.OAuth2.RedirectURL,
				Scopes:       s.config.OAuth2.Scopes,
			}

		})
		core_services_oauth2_github.AddSingletonIGithubOAuth2Authenticator(builder)

	default:
		panic("auth provider not supported")
	}

	switch s.config.AuthStore {
	case "session":
		services_auth_session_token_store.AddScopedITokenStore(builder)
	default:
		services_auth_cookie_token_store.AddScopedITokenStore(builder) // overrides the session one
	}

}

func (s *Startup) addAppHandlers(builder *di.Builder) {

	services_handlers_api_webhook.AddScopedIHandler(builder)
	services_handlers_api_webhookbasicauth.AddScopedIHandler(builder)
	services_handlers_api_webhooknoauth.AddScopedIHandler(builder)
	services_handlers_api_webhookapikey.AddScopedIHandler(builder)

	services_handlers_channel.AddScopedIHandler(builder)

	services_handlers_healthz.AddScopedIHandler(builder)
	services_handlers_ready.AddScopedIHandler(builder)
	services_probes_database.AddSingletonIProbe(builder)
	services_probes_oidc.AddSingletonIProbe(builder)

	services_handlers_home.AddScopedIHandler(builder)
	services_handlers_error.AddScopedIHandler(builder)
	services_handlers_about.AddScopedIHandler(builder)

}

func (s *Startup) ConfigureServices(builder *di.Builder) error {
	fmt.Println(core_utils.PrettyJSON(s.config))
	// add our config as a sigleton object
	di.AddSingletonTypeByObj(builder, s.config)

	services_sse.AddSingletonIServerSideEventServer(builder)
	// Add our main session accessor func
	core_contracts_session.AddGetSessionFunc(builder, app_session.GetSession)
	core_contracts_session.AddGetSessionStoreFunc(builder, s.getSessionStore)

	// Add our secure cookie configs
	s.addSecureCookieOptions(builder)

	// add our auth services
	s.addAuthServices(builder)

	// add our app handlers
	s.addAppHandlers(builder)

	services_claimsprovider.AddSingletonIClaimsProviderMock(builder, s.ctrl)
	services_stores_basicauth.AddMockSingletonIBasicAuthStore(builder, s.ctrl)
	return nil
}
func (s *Startup) Configure(e *echo.Echo, root di.Container) error {
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "header:X-Csrf-Token,form:csrf",
		CookiePath:     "/",
		CookieSecure:   false,
		CookieHTTPOnly: false,
		CookieSameSite: http.SameSiteStrictMode,
		Skipper: func(c echo.Context) bool {
			if core_echo.HasWellknownAuthHeaders(c) {
				return true
			}
			path := c.Request().URL.Path
			res := strings.HasPrefix(path, "/api/")
			return res
		},
	}))

	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			id := uuid.New()
			return id.String()
		},
	}))
	// DevelopmentMiddlewareUsingClaimsMap adds all the needed claims so that FinalAuthVerificationMiddlewareUsingClaimsMap succeeds
	//e.Use(middleware_claimsprincipal.DevelopmentMiddlewareUsingClaimsMap(echostarter_auth.BuildGrpcEntrypointPermissionsClaimsMap(), true))
	e.Use(echo_middleware.JWT(s.GetContainer()))
	e.Use(echo_middleware.BasicAuthWithIBasicAuthStore(s.GetContainer()))
	e.Use(echo_middleware.KeyAuthWithIBasicAuthStore(s.GetContainer()))

	//e.Use(middleware_claimsprincipal.AuthenticatedSessionToClaimsPrincipalMiddleware(root))
	e.Use(core_middleware_claimsprincipal.FinalAuthVerificationMiddlewareUsingClaimsMap(echostarter_auth.BuildGrpcEntrypointPermissionsClaimsMap(), true))
	// only after we pass auth do we slide out the auth session
	e.Use(core_middleware_session.EnsureSlidingSession(root, app_session.GetAuthSession))

	return nil
}
