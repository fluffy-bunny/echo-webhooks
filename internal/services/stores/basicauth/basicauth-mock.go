package basicauth

import (
	"reflect"

	"crypto/subtle"
	contracts_stores_basicauth "echo-starter/internal/contracts/stores/basicauth"
	mocks_stores_basicauth "echo-starter/internal/mocks/stores/basicauth"

	contracts_logger "github.com/fluffy-bunny/grpcdotnetgo/pkg/contracts/logger"
	di "github.com/fluffy-bunny/sarulabsdi"
	"github.com/golang/mock/gomock"
)

type (
	service struct {
		Logger contracts_logger.ILogger `inject:""`
	}
)

func assertImplementation() {
	var _ contracts_stores_basicauth.IBasicAuthStore = (*service)(nil)
	var _ contracts_stores_basicauth.IBasicAuthStore = (*mocks_stores_basicauth.MockIBasicAuthStore)(nil)
}

var reflectType = reflect.TypeOf((*service)(nil))

// AddSingletonIBasicAuthStore registers the *service as a singleton.
func AddSingletonIBasicAuthStore(builder *di.Builder) {
	contracts_stores_basicauth.AddSingletonIBasicAuthStore(builder, reflectType)
}

func AddMockSingletonIBasicAuthStore(builder *di.Builder, mockController *gomock.Controller) {
	mock := mocks_stores_basicauth.NewMockIBasicAuthStore(mockController)
	mock.EXPECT().Validate(gomock.Any(), gomock.Any()).DoAndReturn(func(username string, password string) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
			return true, nil
		}
		return false, nil
	}).AnyTimes()
	contracts_stores_basicauth.AddSingletonIBasicAuthStoreByObj(builder, mock)
}

func (s *service) Validate(username string, password string) (bool, error) {
	return true, nil
}
