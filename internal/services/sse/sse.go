package sse

import (
	"reflect"

	contracts_sse "echo-starter/internal/contracts/sse"

	go_sse "github.com/alexandrevicenzi/go-sse"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	di "github.com/fluffy-bunny/sarulabsdi"
)

type (
	service struct {
		Logger contracts_logger.ILogger `inject:""`
	}
)

func assertImplementation() {
	var _ contracts_sse.IServerSideEventServer = (*go_sse.Server)(nil)
}

var reflectType = reflect.TypeOf((*go_sse.Server)(nil))

// AddSingletonIServerSideEventServer registers the *service as a singleton.
func AddSingletonIServerSideEventServer(builder *di.Builder) {
	contracts_sse.AddSingletonIServerSideEventServerByFunc(builder, reflectType, func(ctn di.Container) (interface{}, error) {
		s := go_sse.NewServer(nil)
		return s, nil
	})
}
