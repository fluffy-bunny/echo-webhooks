package basicauth

//go:generate genny -pkg $GOPACKAGE -in=../../../../genny/sarulabsdi/interface-types.go -out=gen-$GOFILE gen "InterfaceType=IBasicAuthStore"

//go:generate mockgen -package=$GOPACKAGE -destination=../../../mocks/stores/$GOPACKAGE/mock_$GOFILE   echo-starter/internal/contracts/stores/$GOPACKAGE IBasicAuthStore

type (
	// IBasicAuthStore ...
	IBasicAuthStore interface {
		Validate(username string, password string) (bool, error)
	}
)
