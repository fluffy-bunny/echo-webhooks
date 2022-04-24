package channel

import (
	"echo-starter/internal/wellknown"
	"reflect"

	contracts_sse "echo-starter/internal/contracts/sse"

	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"

	contracts_core_claimsprincipal "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/claimsprincipal"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger                contracts_logger.ILogger                        `inject:""`
		ClaimsPrincipal       contracts_core_claimsprincipal.IClaimsPrincipal `inject:""`
		ServerSideEventServer contracts_sse.IServerSideEventServer            `inject:""`
	}
)

func assertImplementation() {
	var _ contracts_handler.IHandler = (*service)(nil)
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddScopedIHandler registers the *service as a singleton.
func AddScopedIHandler(builder *di.Builder) {
	contracts_handler.AddScopedIHandlerEx(builder,
		reflectType,
		[]contracts_handler.HTTPVERB{
			contracts_handler.GET,
			contracts_handler.POST,
			contracts_handler.PUT,
			contracts_handler.DELETE,
			contracts_handler.PATCH,
			contracts_handler.HEAD,
			contracts_handler.OPTIONS,
			contracts_handler.CONNECT,
			contracts_handler.TRACE,
		},
		wellknown.ChannelPath)
}

func (s *service) Ctor() {}
func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {
	req := c.Request()
	res := c.Response()
	s.ServerSideEventServer.ServeHTTP(res, req)
	return nil
}
