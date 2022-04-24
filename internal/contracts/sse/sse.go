package sse

//go:generate genny -pkg $GOPACKAGE -in=../../../genny/sarulabsdi/interface-types.go -out=gen-$GOFILE gen "InterfaceType=IServerSideEventServer"

//go:generate mockgen -package=$GOPACKAGE -destination=../../mocks/$GOPACKAGE/mock_$GOFILE   echo-starter/internal/contracts/$GOPACKAGE IServerSideEventServer
import (
	"net/http"

	"github.com/alexandrevicenzi/go-sse"
)

type (
	// IServerSideEventServer ...
	IServerSideEventServer interface {
		SendMessage(channelName string, message *sse.Message)
		Restart()
		Shutdown()
		ClientCount() int
		HasChannel(name string) bool
		GetChannel(name string) (*sse.Channel, bool)
		Channels() []string
		CloseChannel(name string)
		ServeHTTP(response http.ResponseWriter, request *http.Request)
	}
)
