package wellknown

const (
	HomePath    = "/"
	AboutPath   = "/about"
	HealthzPath = "/healthz"
	ReadyPath   = "/ready"

	WebHookPath          = "/api/:version/webhook"
	WebHookBasicAuthPath = "/api/:version/webhook-basic-auth"
	WebHookNoAuthPath    = "/api/:version/webhook-no-auth"
	WebHookApiKeyPath    = "/api/:version/webhook-api-key"

	ChannelPath = "/events/:channel"
)
