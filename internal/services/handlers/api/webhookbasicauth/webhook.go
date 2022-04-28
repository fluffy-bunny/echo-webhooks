package webhookbasicauth

import (
	contracts_sse "echo-starter/internal/contracts/sse"
	services_handlers_api "echo-starter/internal/services/handlers/api"
	"echo-starter/internal/wellknown"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	go_sse "github.com/alexandrevicenzi/go-sse"
	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	contracts_handler "github.com/fluffy-bunny/grpcdotnetgo/pkg/echo/contracts/handler"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/labstack/echo/v4"
)

type (
	service struct {
		Logger                contracts_logger.ILogger             `inject:""`
		ServerSideEventServer contracts_sse.IServerSideEventServer `inject:""`
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
			contracts_handler.POST,
		},
		wellknown.WebHookBasicAuthPath)
}

func (s *service) GetMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{}
}

func (s *service) Do(c echo.Context) error {

	switch c.Request().Method {
	case http.MethodPost:
		return s.post(c)
	default:
		return echo.NewHTTPError(http.StatusMethodNotAllowed)
	}
}

type echoRequest struct {
	Path   string      `json:"path"`
	Method string      `json:"method"`
	Header http.Header `json:"header"`
	Body   interface{} `json:"body"`
}

func (s *service) post(c echo.Context) error {

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to read request body")
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read request body")
	}
	body := make(map[string]interface{})
	if len(b) > 0 {
		err = json.Unmarshal(b, &body)
		if err != nil {
			s.Logger.Error().Err(err).Msg("failed to read request body")
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to read request body")
		}
	}

	response := echoRequest{
		Path:   c.Request().URL.Path,
		Method: c.Request().Method,
		Header: c.Request().Header,
		Body:   body,
	}
	b, err = json.Marshal(response)
	s.ServerSideEventServer.SendMessage("/events/webhooks", go_sse.SimpleMessage(string(b)))

	services_handlers_api.WebhookCount.Inc()
	count := fmt.Sprintf("%d", services_handlers_api.WebhookCount.Load())
	s.ServerSideEventServer.SendMessage("/events/webhooks-progress", go_sse.SimpleMessage(count))

	return c.JSON(http.StatusOK, response)
}
